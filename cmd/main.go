package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/CleverseAcademy/cd-compose-deployment/api"
	"github.com/CleverseAcademy/cd-compose-deployment/api/auth"
	"github.com/CleverseAcademy/cd-compose-deployment/api/services"
	"github.com/CleverseAcademy/cd-compose-deployment/config"
	"github.com/CleverseAcademy/cd-compose-deployment/constants"
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/CleverseAcademy/cd-compose-deployment/providers"
	"github.com/CleverseAcademy/cd-compose-deployment/usecases"
	"github.com/docker/docker/client"
	"github.com/gofiber/fiber/v2"
)

func main() {
	clnt, err := client.NewClientWithOpts(client.WithHostFromEnv(), client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	prj, err := providers.LoadComposeProject(providers.IArgsLoadComposeProject{
		WorkingDir:  config.AppConfig.ComposeWorkingDir,
		ComposeFile: config.AppConfig.ComposeFile,
		ProjectName: config.AppConfig.ComposeProjectName,
	})
	if err != nil {
		panic(err)
	}

	walWriter, err := providers.CreateWalWriter(filepath.Join(config.AppConfig.DataDir, constants.EventLogDir))
	if err != nil {
		panic(err)
	}

	eventLogBase := usecases.EventLogUseCase{
		Logger:       walWriter,
		DockerClient: clnt,
	}

	useCaseLogDeploymentDoneEvent := &usecases.UseCaseLogDeploymentDoneEvent{
		EventLogUseCase: &eventLogBase,
	}
	useCaseLogDeploymentFailureEvent := &usecases.UseCaseLogDeploymentFailureEvent{
		EventLogUseCase: &eventLogBase,
	}
	useCaseLogConfigLoadedEvent := &usecases.UseCaseLogConfigLoadedEvent{
		EventLogUseCase: &eventLogBase,
	}
	useCaseLogDeplomentSkipped := &usecases.UseCaseLogDeploymentSkippedEvent{
		EventLogUseCase: &eventLogBase,
	}

	err = useCaseLogConfigLoadedEvent.Execute(*prj)
	if err != nil {
		panic(err)
	}

	deploymentBase := usecases.DeploymentUseCase{
		Project: *prj,
	}

	useCasePrepareServiceDeployment := &usecases.UseCasePrepareServiceDeployment{
		DeploymentUseCase: &deploymentBase,
	}
	useCaseEnqueueServiceDeployment := &usecases.UseCaseEnqueueServiceDeployment{
		DeploymentUseCase: &deploymentBase,
		Logs:              &entities.DeploymentTable{},
	}
	useCaseExecuteServiceDeployments := &usecases.UseCaseExecuteServiceDeployments{
		UseCaseEnqueueServiceDeployment: useCaseEnqueueServiceDeployment,
	}
	useCaseGetAllServiceDeploymentInfo := &usecases.UseCaseGetAllServiceDeploymentInfo{
		UseCaseEnqueueServiceDeployment: useCaseEnqueueServiceDeployment,
	}

	service := services.Service{
		GetAllServiceDeploymentInfo: useCaseGetAllServiceDeploymentInfo,
		ExecuteServiceDeployments:   useCaseExecuteServiceDeployments,
		LogDeploymentDoneEvent:      useCaseLogDeploymentDoneEvent,
		LogDeploymentFailureEvent:   useCaseLogDeploymentFailureEvent,
		LogDeploymentSkippedEvent:   useCaseLogDeplomentSkipped,
	}

	composeAPI, err := providers.GetComposeService(clnt, config.AppConfig.DockerContext)
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	authMDW := auth.SignatureVerificationMiddleware(auth.IArgsCreateSignatureVerificationMiddleware{
		IService: service,
	})

	app.Post(
		constants.PathAddDeployment,
		authMDW,
		api.DeployNewImageHandler(api.IArgsCreateDeployNewImageHandler{
			PrepareServiceDeployment: useCasePrepareServiceDeployment,
			EnqueueServiceDeployment: useCaseEnqueueServiceDeployment,
		}))

	app.Get(
		constants.PathGetDeploymentJTI,
		api.GetNextDeploymentJTIHandler(api.IArgsCreateGetNextDeploymentJTIHandler{
			IService: service,
		}))

	go func(s services.Service) {
		for {
			time.Sleep(config.AppConfig.DeployInterval)

			for _, svc := range prj.Services {
				_, err := s.SoyDeploy(services.IArgsCreateDeployNewImageHandler{
					ServiceName: svc.Name,
					ComposeAPI:  composeAPI,
					DockerClnt:  clnt,
				})
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}(service)

	err = app.Listen(config.AppConfig.ListeningSocket)
	if err != nil {
		panic(err)
	}
}
