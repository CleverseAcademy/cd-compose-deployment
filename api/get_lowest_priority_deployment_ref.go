package api

import (
	"strings"

	"github.com/CleverseAcademy/cd-compose-deployment/constants"
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/CleverseAcademy/cd-compose-deployment/usecases"
	"github.com/gofiber/fiber/v2"
)

type IArgsGetLowestPriorityDeploymentRef struct {
	GetLowestPriorityDeploymentInfo usecases.IUseCaseGetLowestPriorityDeploymentInfo
}

func GetDeploymentRef(args IArgsGetLowestPriorityDeploymentRef) fiber.Handler {
	return func(c *fiber.Ctx) error {
		serviceName := c.Params("serviceName")

		d, err := args.GetLowestPriorityDeploymentInfo.Execute(entities.ServiceName(serviceName))
		if err != nil {
			if strings.Contains(err.Error(), constants.ErrorEmptyDeployment) {
				return c.SendStatus(fiber.StatusNoContent)
			}

			return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
		}

		return c.Status(fiber.StatusOK).JSON(d)
	}
}
