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

type UseCaseLogConfigLoadedEvent struct {
	*EventLogUseCase
}

func (u *UseCaseLogConfigLoadedEvent) Execute(prj types.Project) error {
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

	servicesInfo := make([]entities.ServiceInfo, len(services))

	for idx, service := range services {
		servicesInfo[idx] = entities.ServiceInfo{
			Name:  service.Name,
			Image: service.Image,
		}
		if service.Container != nil {
			servicesInfo[idx].ContainerID = service.Container.ID
		}
	}

	event := entities.ConfigLoadedEventEntry{
		Services: servicesInfo,
		EventLog: entities.EventLog{
			Timestamp:       t,
			Event:           constants.ConfigLoadedEventName,
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
