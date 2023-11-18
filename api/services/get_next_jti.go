package services

import (
	"strings"

	"github.com/CleverseAcademy/cd-compose-deployment/config"
	"github.com/CleverseAcademy/cd-compose-deployment/constants"
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/CleverseAcademy/cd-compose-deployment/utils"
	"github.com/pkg/errors"
)

func (s Service) GetNextJTI(serviceName string) (string, error) {
	deployments, err := s.GetAllServiceDeploymentInfo.Execute(entities.ServiceName(serviceName))
	if err != nil && strings.Contains(err.Error(), constants.ErrorServiceNotFound) {
		return "", errors.New(constants.ErrorServiceNotFound)
	}

	nextJti, err := utils.Base64EncodedSha256([]interface{}{config.AppConfig.InitialHash, deployments})
	return nextJti, errors.Wrap(err, "GetNextJTI@utils.Base64EncodedSha256")
}
