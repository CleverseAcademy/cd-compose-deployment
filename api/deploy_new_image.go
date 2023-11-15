package api

import (
	"fmt"

	"github.com/CleverseAcademy/cd-compose-deployment/api/dto"
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/CleverseAcademy/cd-compose-deployment/usecases"
	"github.com/gofiber/fiber/v2"
)

type IArgsCreateDeployNewImageHandler struct {
	PrepareServiceDeployment usecases.IUseCasePrepareServiceDeployment
	EnqueueServiceDeployment usecases.IUseCaseEnqueueServiceDeployment
}

func DeployNewImageHandler(args IArgsCreateDeployNewImageHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := new(dto.DeployImageDto)

		err := c.BodyParser(request)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		serviceName := entities.ServiceName(request.Service)
		deployment, err := args.PrepareServiceDeployment.Execute(
			serviceName,
			request.Priority,
			request.Ref,
			request.Image,
		)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		currentDeployments := args.EnqueueServiceDeployment.Execute(serviceName, deployment)

		return c.Status(fiber.StatusAccepted).SendString(fmt.Sprintf("%d", currentDeployments))
	}
}
