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

type UseCaseLogStopSignalReceivedEvent struct {
	*EventLogUseCase
}

func (u *UseCaseLogStopSignalReceivedEvent) Execute(prj types.Project) error {
	t := time.Now()
	projectChecksum, err := utils.Base64EncodedSha256(prj)
	if err != nil {
		return errors.Wrap(err, "UseCaseLogStopSignalReceivedEvent")
	}

	event := entities.StopSignalReceivedEventEntry{
		EventLog: entities.EventLog{
			Event:           constants.StopSignalReceivedEventName,
			Timestamp:       t,
			ProjectChecksum: projectChecksum,
		},
	}

	data, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "UseCaseLogStopSignalReceivedEvent@Marshal")
	}

	_, err = u.Logger.Write(data)
	return errors.Wrap(err, "UseCaseLogStopSignalReceivedEvent@Logger.Write")
}
