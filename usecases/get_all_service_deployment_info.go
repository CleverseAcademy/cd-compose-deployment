package usecases

import (
	"github.com/CleverseAcademy/cd-compose-deployment/constants"
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/pkg/errors"
)

type UseCaseGetAllServiceDeploymentInfo struct {
	*UseCaseEnqueueServiceDeployment
}

func (u *UseCaseGetAllServiceDeploymentInfo) Execute(service entities.ServiceName) ([]entities.Deployment, error) {
	u.RLock()
	defer u.RUnlock()

	queue, err := u.Logs.GetServiceDeploymentQueue(service)
	if err != nil {
		for _, svc := range u.Project.Services {
			if svc.Name == string(service) {
				return []entities.Deployment{}, errors.Wrap(err, "GetServiceDeploymentQueue")
			}
		}

		return nil, errors.Wrap(err, constants.ErrorServiceNotFound)
	}

	return queue.Items(), nil
}
