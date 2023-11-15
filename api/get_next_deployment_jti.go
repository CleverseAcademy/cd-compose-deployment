package api

import (
	"github.com/CleverseAcademy/cd-compose-deployment/config"
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/CleverseAcademy/cd-compose-deployment/usecases"
	"github.com/CleverseAcademy/cd-compose-deployment/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type IArgsCreateGetNextDeploymentJTIHandler struct {
	GetAllServiceDeploymentInfo usecases.IUseCaseGetAllServiceDeploymentInfo
}

func GetNextDeploymentJTIHandler(args IArgsCreateGetNextDeploymentJTIHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		deployments, _ := args.GetAllServiceDeploymentInfo.Execute(entities.ServiceName(c.Params("serviceName")))

		nextJti, err := utils.Base64EncodedSha256([]interface{}{config.AppConfig.InitialHash, deployments})
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "GetNextDeploymentJTIHandler").Error())
		}
		return c.SendString(nextJti)
	}
}
