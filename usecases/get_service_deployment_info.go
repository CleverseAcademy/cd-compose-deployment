package usecases

import (
	"fmt"

	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/pkg/errors"
)

type UseCaseGetServiceDeploymentInfo struct {
	*UseCaseEnqueueServiceDeployment
}

func (u *UseCaseGetServiceDeploymentInfo) Execute(service entities.ServiceName, ref string) (entities.Deployment, error) {
	u.RLock()
	defer u.RUnlock()

	queue, err := u.Logs.GetServiceDeploymentQueue(service)
	if err != nil {
		return entities.Deployment{}, errors.Wrap(err, "GetServiceDeploymentQueue")
	}

	deployments := queue.Items()
	for _, deployment := range deployments {
		if deployment.Ref == ref {
			return deployment, nil
		}
	}

	return entities.Deployment{}, fmt.Errorf("ref %s not found", ref)
}
