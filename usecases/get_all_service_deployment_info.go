package usecases

import (
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/pkg/errors"
)

type UseCaseGetAllServiceDeploymentInfo struct {
	*UseCaseEnqueueServiceDeployment
}

func (u *UseCaseGetAllServiceDeploymentInfo) Execute(service entities.ServiceName) ([]entities.Deployment, error) {
	queue, err := u.tbl.GetServiceDeploymentQueue(service)
	if err != nil {
		return []entities.Deployment{}, errors.Wrap(err, "service not found")
	}

	return queue.Items(), nil
}
