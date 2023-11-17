package usecases

import (
	"encoding/json"
	"time"

	"github.com/CleverseAcademy/cd-compose-deployment/constants"
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/CleverseAcademy/cd-compose-deployment/utils"
	"github.com/compose-spec/compose-go/types"
	"github.com/pkg/errors"
)

type UseCaseLogDeploymentFailureEvent struct {
	*EventLogUseCase
}

func (u *UseCaseLogDeploymentFailureEvent) Execute(prj types.Project, failedService entities.UndeployableServiceInfo) error {
	t := time.Now()
	projectChecksum, err := utils.Base64EncodedSha256(prj)
	if err != nil {
		return errors.Wrap(err, "UseCaseLogDeploymentFailureEvent")
	}

	containers, err := utils.GetProjectContainers(prj.Name, u.DockerClient)
	if err != nil {
		return errors.Wrap(err, "UseCaseLogDeploymentFailureEvent")
	}

	filteredServices := make([]types.ServiceConfig, 0)
	for _, svc := range prj.Services {
		if svc.Name == string(failedService.Name) {
			continue
		}

		filteredServices = append(filteredServices, svc)
	}

	services, err := mapServiceStatus(t, filteredServices, containers)
	if err != nil {
		return err
	}

	event := &entities.DeploymentFailureEventEntry{
		EventLog: entities.EventLog{
			Timestamp:       t,
			Event:           constants.DeploymentFailureEventName,
			ProjectChecksum: projectChecksum,
		},
		Services:      services,
		FailedService: failedService,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "UseCaseLogDeploymentFailureEvent@Marshal")
	}

	_, err = u.Logger.Write(data)
	return errors.Wrap(err, "UseCaseLogDeploymentFailureEvent@Logger.Write")
}
