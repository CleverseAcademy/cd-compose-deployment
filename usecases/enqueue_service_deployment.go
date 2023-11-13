package usecases

import (
	"container/heap"

	"github.com/CleverseAcademy/cd-compose-deployment/entities"
)

type UseCaseEnqueueServiceDeployment struct {
	*DeploymentUseCase
	Tbl *entities.DeploymentTable
}

func (u *UseCaseEnqueueServiceDeployment) Execute(service entities.ServiceName, deployment *entities.Deployment) int8 {
	if u.Tbl == nil {
		u.Tbl = &entities.DeploymentTable{}
	}
	queue, err := u.Tbl.GetServiceDeploymentQueue(service)
	if err != nil {
		queue = u.Tbl.InitializeDeploymentQueue(service)
	}

	heap.Init(queue)

	heap.Push(queue, *deployment)

	return int8(queue.Len())
}
