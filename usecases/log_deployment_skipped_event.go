package usecases

import (
	"encoding/json"
	"time"

	"github.com/CleverseAcademy/cd-compose-deployment/constants"
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/compose-spec/compose-go/types"
	"github.com/pkg/errors"
)

type UseCaseLogDeploymentSkippedEvent struct {
	*EventLogUseCase
}

func (u *UseCaseLogDeploymentSkippedEvent) Execute(prj types.Project, serviceName entities.ServiceName) error {
	t := time.Now()

	event := entities.DeploymentSkippedEventEntry{
		EventLog: entities.EventLog{Event: constants.DeploymentSkippedEventName, Timestamp: t},
		Name:     serviceName,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "UseCaseLogDeploymentSkippedEvent@Marshal")
	}

	_, err = u.Logger.Write(data)
	return errors.Wrap(err, "UseCaseLogDeploymentSkippedEvent@Logger.Write")
}
