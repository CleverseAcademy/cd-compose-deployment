package services

import (
	"github.com/CleverseAcademy/cd-compose-deployment/config"
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/CleverseAcademy/cd-compose-deployment/utils"
	"github.com/pkg/errors"
)

type IArgsGenerateJTI struct {
	ServiceName string
}

func (s Service) GetNextJTI(serviceName string) (string, error) {
	deployments, _ := s.GetAllServiceDeploymentInfo.Execute(entities.ServiceName(serviceName))

	nextJti, err := utils.Base64EncodedSha256([]interface{}{config.AppConfig.InitialHash, deployments})
	return nextJti, errors.Wrap(err, "GenerateJTI")
}
