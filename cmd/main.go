package main

import (
	"path/filepath"

	"github.com/CleverseAcademy/cd-compose-deployment/api"
	"github.com/CleverseAcademy/cd-compose-deployment/api/auth"
	"github.com/CleverseAcademy/cd-compose-deployment/api/logger"
	"github.com/CleverseAcademy/cd-compose-deployment/api/services"
	"github.com/CleverseAcademy/cd-compose-deployment/config"
	"github.com/CleverseAcademy/cd-compose-deployment/constants"
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

	eventsLogger, err := providers.CreateWalWriter(filepath.Join(config.AppConfig.DataDir, constants.EventLogDir))
	if err != nil {
		panic(err)
	}

	accessLogger, err := providers.CreateWalWriter(filepath.Join(config.AppConfig.DataDir, constants.AccessLogDir))
	if err != nil {
		panic(err)
	}

	entropy := providers.NewEntropy([]byte(config.AppConfig.InitialHash))
	err = accessLogger.RegisterEntropyObserver(entropy)
	if err != nil {
		panic(err)
	}

	eventLogBase := usecases.EventLogUseCase{
		Logger:       eventsLogger,
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
	err = eventsLogger.RegisterEntropyObserver(entropy)
	if err != nil {
		panic(err)
	}

	deploymentBase := usecases.DeploymentUseCase{
		Project: *prj,
	}

	useCaseEnqueueServiceDeployment := usecases.CreateUseCaseEnqueueServiceDeployment(&deploymentBase)

	useCasePrepareServiceDeployment := &usecases.UseCasePrepareServiceDeployment{
		DeploymentUseCase: &deploymentBase,
	}
	useCaseExecuteServiceDeployments := &usecases.UseCaseExecuteServiceDeployments{
		UseCaseEnqueueServiceDeployment: useCaseEnqueueServiceDeployment,
	}
	useCaseGetCurrentHighestPriorityDeploymentInfo := &usecases.UseCaseGetCurrentHighestPriorityDeploymentInfo{
		UseCaseEnqueueServiceDeployment: useCaseEnqueueServiceDeployment,
	}

	service := services.Service{
		ExecuteServiceDeployments: useCaseExecuteServiceDeployments,
		LogDeploymentDoneEvent:    useCaseLogDeploymentDoneEvent,
		LogDeploymentFailureEvent: useCaseLogDeploymentFailureEvent,
		LogDeploymentSkippedEvent: useCaseLogDeplomentSkipped,
	}

	composeAPI, err := providers.GetComposeService(clnt, config.AppConfig.DockerContext)
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	getRequestAuthMDW := auth.SignatureVerificationMiddleware(auth.IArgsCreateSignatureVerificationMiddleware{
		Entropy:           entropy,
		VerifyRequestBody: false,
	})
	postRequestAuthMDW := auth.SignatureVerificationMiddleware(auth.IArgsCreateSignatureVerificationMiddleware{
		Entropy:           entropy,
		VerifyRequestBody: true,
	})
	accessLogMDW := logger.NewAccessLogMiddleware(accessLogger)

	app.Post(
		constants.PathAddDeployment,
		postRequestAuthMDW,
		accessLogMDW,
		api.DeployNewImageHandler(api.IArgsCreateDeployNewImageHandler{
			PrepareServiceDeployment: useCasePrepareServiceDeployment,
			EnqueueServiceDeployment: useCaseEnqueueServiceDeployment,
		}))

	app.Get(
		constants.PathGetJTI,
		api.GetNextJTIHandler(api.IArgsCreateGetNextDeploymentJTIHandler{
			Entropy: entropy,
		}))

	app.Get(
		constants.PathGetNextDeploymentInfo,
		getRequestAuthMDW,
		accessLogMDW,
		api.GetNextDeployment(api.IArgsGetCurrentHighestPriorityDeploymentInfo{
			GetCurrentHighestPriorityDeploymentInfo: useCaseGetCurrentHighestPriorityDeploymentInfo,
		}),
	)

	go service.PeriodicallySoyDeploy(clnt, composeAPI, prj, filepath.Join(config.AppConfig.DataDir, constants.DefaultComposeYMLFilename))

	err = app.Listen(config.AppConfig.ListeningSocket)
	if err != nil {
		panic(err)
	}
}
