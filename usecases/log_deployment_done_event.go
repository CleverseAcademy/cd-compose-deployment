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

type UseCaseLogDeploymentDoneEvent struct {
	*EventLogUseCase
}

func (u *UseCaseLogDeploymentDoneEvent) Execute(prj types.Project, deployment entities.Deployment, serviceName entities.ServiceName) error {
	t := time.Now()
	projectChecksum, err := utils.Base64EncodedSha256(prj)
	if err != nil {
		return errors.Wrap(err, "utils.Base64EncodedSha256")
	}

	containers, err := utils.GetProjectContainers(prj.Name, u.DockerClient)
	if err != nil {
		return errors.Wrap(err, "utils.GetProjectContainers")
	}

	services, err := mapServiceStatus(t, prj.Services, containers)
	if err != nil {
		return errors.Wrap(err, "mapServiceStatus")
	}

	for idx, s := range services {
		if serviceName == s.Name {
			services[idx].DeploymentRef = deployment.Ref
			services[idx].ServiceDetail.Timestamp = deployment.Timestamp
		}
	}

	event := &entities.DeploymentDoneEventEntry{
		EventLog: entities.EventLog{
			Timestamp:       t,
			Event:           constants.DeploymentDoneEventName,
			ProjectChecksum: projectChecksum,
		},
		Services: services,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	_, err = u.Logger.Write(data)
	return errors.Wrap(err, "Logger.Write")
}
