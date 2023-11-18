package usecases

import (
	"time"

	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/CleverseAcademy/cd-compose-deployment/utils"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/compose/v2/pkg/api"
	docker "github.com/docker/docker/api/types"
	"github.com/pkg/errors"
)

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
