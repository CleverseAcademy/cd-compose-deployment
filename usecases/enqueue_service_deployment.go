package usecases

import (
	"container/heap"

	"github.com/CleverseAcademy/cd-compose-deployment/entities"
)

type UseCaseEnqueueServiceDeployment struct {
	*DeploymentUseCase
	tbl entities.DeploymentTable
}

func (u UseCaseEnqueueServiceDeployment) Execute(service entities.ServiceName, deployment *entities.Deployment) int8 {
	queue, err := u.tbl.GetServiceDeploymentQueue(service)
	if err != nil {
		queue = u.tbl.InitializeDeploymentQueue(service)
	}

	heap.Init(queue)

	heap.Push(queue, deployment)

	return int8(queue.Len())
}
