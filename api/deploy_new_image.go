package api

import (
	"github.com/CleverseAcademy/cd-compose-deployment/api/dto"
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/CleverseAcademy/cd-compose-deployment/usecases"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/docker/client"
	"github.com/gofiber/fiber/v2"
)

type IArgsCreateDeployNewImageHandler struct {
	ComposeAPI                api.Service
	DockerClnt                *client.Client
	PrepareServiceDeployment  usecases.IUseCasePrepareServiceDeployment
	EnqueueServiceDeployment  usecases.IUseCaseEnqueueServiceDeployment
	ExecuteServiceDeployments usecases.IUseCaseExecuteServiceDeployments
}

func DeployNewImageHandler(args IArgsCreateDeployNewImageHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(dto.DeployImageDto)

		err := c.BodyParser(request)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		targetService := entities.ServiceName(request.Service)
		deployment, err := args.PrepareServiceDeployment.Execute(
			targetService,
			request.Priority,
			request.Ref,
			request.Image,
		)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		currentDeployments := args.EnqueueServiceDeployment.Execute(targetService, deployment)

		prj, err := args.ExecuteServiceDeployments.Execute(args.DockerClnt, args.ComposeAPI, targetService)
		if err != nil {
			return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
		}

		return c.JSON([]interface{}{prj, currentDeployments})
	}
}
