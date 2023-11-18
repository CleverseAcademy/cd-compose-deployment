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

type UseCaseLogConfigChangesDetectedEvent struct {
	*EventLogUseCase
}

func (u *UseCaseLogConfigChangesDetectedEvent) Execute(prj types.Project) error {
	t := time.Now()
	projectChecksum, err := utils.Base64EncodedSha256(prj)
	if err != nil {
		return errors.Wrap(err, "utils.Base64EncodedSha256")
	}

	event := entities.ConfigChangesDetectedEventEntry{
		EventLog: entities.EventLog{
			Event:           constants.ConfigChangesDetectedEventName,
			Timestamp:       t,
			ProjectChecksum: projectChecksum,
		},
	}

	data, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	_, err = u.Logger.Write(data)
	return errors.Wrap(err, "Logger.Write")
}
