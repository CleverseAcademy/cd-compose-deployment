package usecases

import (
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/pkg/errors"
)

type UseCaseGetLatestServiceDeploymentInfo struct {
	*UseCaseEnqueueServiceDeployment
}

func (u *UseCaseGetLatestServiceDeploymentInfo) Execute(service entities.ServiceName) (entities.Deployment, error) {
	queue, err := u.Tbl.GetServiceDeploymentQueue(service)
	if err != nil {
		return entities.Deployment{}, errors.Wrap(err, "service not found")
	}

	deployments := queue.Items()
	latestDeployment := deployments[len(deployments)-1]

	for _, svc := range deployments {
		if svc.Timestamp.After(latestDeployment.Timestamp) {
			latestDeployment = svc
		}
	}

	return latestDeployment, nil
}
