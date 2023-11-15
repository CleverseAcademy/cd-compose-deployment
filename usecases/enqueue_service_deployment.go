package usecases

import (
	"container/heap"

	"github.com/CleverseAcademy/cd-compose-deployment/entities"
)

type UseCaseEnqueueServiceDeployment struct {
	*DeploymentUseCase
	Logs *entities.DeploymentTable
	tbl  *entities.DeploymentTable
}

func (u *UseCaseEnqueueServiceDeployment) Execute(service entities.ServiceName, deployment *entities.Deployment) int8 {
	u.Lock()
	defer u.Unlock()

	if u.tbl == nil {
		u.tbl = &entities.DeploymentTable{}
	}
	queue, err := u.tbl.GetServiceDeploymentQueue(service)
	if err != nil {
		queue = u.tbl.InitializeDeploymentQueue(service)
	}

	logs, err := u.Logs.GetServiceDeploymentQueue(service)
	if err != nil {
		logs = u.Logs.InitializeDeploymentQueue(service)
	}

	logs.Push(*deployment)

	heap.Init(queue)

	heap.Push(queue, *deployment)

	return int8(queue.Len())
}
