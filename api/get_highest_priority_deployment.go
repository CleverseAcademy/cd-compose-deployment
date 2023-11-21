package api

import (
	"strings"

	"github.com/CleverseAcademy/cd-compose-deployment/constants"
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/CleverseAcademy/cd-compose-deployment/usecases"
	"github.com/gofiber/fiber/v2"
)

type IArgsGetCurrentHighestPriorityDeploymentInfo struct {
	GetCurrentHighestPriorityDeploymentInfo usecases.IUseCaseGetCurrentHighestPriorityDeploymentInfo
}

func GetNextDeployment(args IArgsGetCurrentHighestPriorityDeploymentInfo) fiber.Handler {
	return func(c *fiber.Ctx) error {
		serviceName := c.Params("serviceName")

		d, err := args.GetCurrentHighestPriorityDeploymentInfo.Execute(entities.ServiceName(serviceName))
		if err != nil {
			if strings.Contains(err.Error(), constants.ErrorEmptyDeployment) {
				return c.SendStatus(fiber.StatusNoContent)
			}

			return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
		}

		return c.Status(fiber.StatusOK).JSON(d)
	}
}
