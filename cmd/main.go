package main

import (
	"github.com/CleverseAcademy/cd-compose-deployment/api"
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

	workingDir := "/Users/intaniger/works/focusing/CleverseAcademy/learnhub-api"

	composeAPI, err := providers.GetComposeService(clnt, "desktop-linux")
	if err != nil {
		panic(err)
	}

	prj, err := providers.LoadComposeProject(providers.IArgsLoadComposeProject{
		WorkingDir: workingDir,
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
	}

	useCaseExecuteServiceDeployments := &usecases.UseCaseExecuteServiceDeployments{
		UseCaseEnqueueServiceDeployment: useCaseEnqueueServiceDeployment,
	}

	app := fiber.New()

	app.Post("/deploy", api.DeployNewImageHandler(api.IArgsCreateDeployNewImageHandler{
		ComposeAPI:                composeAPI,
		PrepareServiceDeployment:  useCasePrepareServiceDeployment,
		EnqueueServiceDeployment:  useCaseEnqueueServiceDeployment,
		ExecuteServiceDeployments: useCaseExecuteServiceDeployments,
	}))

	err = app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
