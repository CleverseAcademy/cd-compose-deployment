package usecases

import (
	"encoding/json"
	"time"

	"github.com/CleverseAcademy/cd-compose-deployment/constants"
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/pkg/errors"
)

type UseCaseLogDeploymentSkippedEvent struct {
	*EventLogUseCase
}

func (u *UseCaseLogDeploymentSkippedEvent) Execute(serviceName entities.ServiceName) error {
	t := time.Now()

	event := entities.DeploymentSkippedEventEntry{
		EventLog: entities.EventLog{Event: constants.DeploymentSkippedEventName, Timestamp: t},
		Name:     serviceName,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	_, err = u.Logger.Write(data)
	return errors.Wrap(err, "Logger.Write")
}
