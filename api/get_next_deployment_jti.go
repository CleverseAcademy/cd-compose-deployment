package api

import (
	"github.com/CleverseAcademy/cd-compose-deployment/api/services"
	"github.com/gofiber/fiber/v2"
)

type IArgsCreateGetNextDeploymentJTIHandler struct {
	services.IService
}

func GetNextDeploymentJTIHandler(args IArgsCreateGetNextDeploymentJTIHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		nextJti, err := args.GetNextJTI(c.Params("serviceName"))
		if err != nil {
			return fiber.NewError(fiber.StatusServiceUnavailable, err.Error())
		}

		return c.SendString(nextJti)
	}
}
