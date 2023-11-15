package main

import (
	"github.com/CleverseAcademy/cd-compose-deployment/api"
	"github.com/CleverseAcademy/cd-compose-deployment/api/auth"
	"github.com/CleverseAcademy/cd-compose-deployment/api/services"
	"github.com/CleverseAcademy/cd-compose-deployment/config"
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

	composeAPI, err := providers.GetComposeService(clnt, config.AppConfig.DockerContext)
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

	base := usecases.DeploymentUseCase{
		Project: *prj,
	}

	useCasePrepareServiceDeployment := &usecases.UseCasePrepareServiceDeployment{
		DeploymentUseCase: &base,
	}

	useCaseEnqueueServiceDeployment := &usecases.UseCaseEnqueueServiceDeployment{
		DeploymentUseCase: &base,
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
	}

	app := fiber.New()

	authMDW := auth.SignatureVerificationMiddleware(auth.IArgsCreateSignatureVerificationMiddleware{
		IService: service,
	})

	app.Post(
		"/deploy",
		authMDW,
		api.DeployNewImageHandler(api.IArgsCreateDeployNewImageHandler{
			DockerClnt:                clnt,
			ComposeAPI:                composeAPI,
			PrepareServiceDeployment:  useCasePrepareServiceDeployment,
			EnqueueServiceDeployment:  useCaseEnqueueServiceDeployment,
			ExecuteServiceDeployments: useCaseExecuteServiceDeployments,
		}))

	app.Get(
		"/deploy/nextJTI/:serviceName",
		api.GetNextDeploymentJTIHandler(api.IArgsCreateGetNextDeploymentJTIHandler{
			IService: service,
		}))

	err = app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
