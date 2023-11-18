package api

import (
	"github.com/CleverseAcademy/cd-compose-deployment/providers"
	"github.com/gofiber/fiber/v2"
)

type IArgsCreateGetNextDeploymentJTIHandler struct {
	*providers.Entropy
}

func GetNextJTIHandler(args IArgsCreateGetNextDeploymentJTIHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString(args.Base64Get())
	}
}
