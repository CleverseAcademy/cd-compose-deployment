package api

import (
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/CleverseAcademy/cd-compose-deployment/usecases"
	"github.com/CleverseAcademy/cd-compose-deployment/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type IArgsCreateListAllDeploymentsHandler struct {
	GetAllServiceDeploymentInfo usecases.IUseCaseGetAllServiceDeploymentInfo
}

func ListAllDeploymentsHandler(args IArgsCreateListAllDeploymentsHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		deployments, err := args.GetAllServiceDeploymentInfo.Execute(entities.ServiceName(c.Params("serviceName")))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "ListAllDeploymentsHandler").Error())
		}

		jti, err := utils.Base64EncodedSha256(deployments)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "ListAllDeploymentsHandler").Error())
		}
		return c.SendString(jti)
	}
}
