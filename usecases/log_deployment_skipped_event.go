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

type UseCaseLogDeploymentSkippedEvent struct {
	*EventLogUseCase
}

func (u *UseCaseLogDeploymentSkippedEvent) Execute(prj types.Project, serviceName entities.ServiceName) error {
	t := time.Now()
	projectChecksum, err := utils.Base64EncodedSha256(prj)
	if err != nil {
		return errors.Wrap(err, "UseCaseLogDeploymentSkippedEvent")
	}

	event := entities.DeploymentSkippedEventEntry{
		EventLog: entities.EventLog{Event: constants.DeploymentSkippedEventName, Timestamp: t, ProjectChecksum: projectChecksum},
		Name:     serviceName,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "UseCaseLogDeploymentSkippedEvent@Marshal")
	}

	_, err = u.Logger.Write(data)
	return errors.Wrap(err, "UseCaseLogDeploymentSkippedEvent@Logger.Write")
}
