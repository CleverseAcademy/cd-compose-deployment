package usecases

import (
	"encoding/json"
	"time"

	"github.com/CleverseAcademy/cd-compose-deployment/constants"
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/CleverseAcademy/cd-compose-deployment/utils"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/compose/v2/pkg/api"
	docker "github.com/docker/docker/api/types"
	"github.com/pkg/errors"
)

type UseCaseLogDeploymentDoneEvent struct {
	*EventLogUseCase
}

func (u *UseCaseLogDeploymentDoneEvent) Execute(prj types.Project, deployment entities.Deployment, deployedServiceName entities.ServiceName) error {
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
		return err
	}

	for idx, s := range services {
		if deployedServiceName == s.Name {
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

func mapServiceStatus(t time.Time, services types.Services, containers []docker.Container) ([]entities.DetailedServiceInfo, error) {
	containerServiceDict := make(map[string]*docker.Container)
	for idx, c := range containers {
		containerServiceDict[c.Labels[api.ServiceLabel]] = &containers[idx]
	}

	svcStatus := make([]entities.DetailedServiceInfo, len(services))
	for idx, svc := range services {
		cfg := svc
		cfg.Image = ""

		checksum, err := utils.Base64EncodedSha256(cfg)
		if err != nil {
			return nil, errors.Wrap(err, "ServicesLoop.Base64EncodedSha256")
		}

		svcStatus[idx] = entities.DetailedServiceInfo{
			Name:          entities.ServiceName(svc.Name),
			DeploymentRef: "",
			Container:     nil,
			ServiceDetail: &entities.ServiceDetail{Timestamp: t, CfgChecksum: checksum},
		}

		container, ok := containerServiceDict[svc.Name]
		if ok {
			svcStatus[idx].Container = container
			svcStatus[idx].ServiceDetail.Image = container.ID
		}
	}
	return svcStatus, nil
}
