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

func CreateUseCaseEnqueueServiceDeployment(
	base *DeploymentUseCase,
) *UseCaseEnqueueServiceDeployment {
	return &UseCaseEnqueueServiceDeployment{
		DeploymentUseCase: base,
		Logs:              entities.NewDeploymentTable(),
		tbl:               entities.NewDeploymentTable(),
	}
}

func (u *UseCaseEnqueueServiceDeployment) Execute(service entities.ServiceName, deployment *entities.Deployment) int8 {
	u.Lock()
	defer u.Unlock()

	queue, err := u.tbl.GetServiceDeploymentQueue(service)
	if err != nil {
		queue = u.tbl.InitializeDeploymentQueue(service)
	}

	logs, err := u.Logs.GetServiceDeploymentQueue(service)
	if err != nil {
		logs = u.Logs.InitializeDeploymentQueue(service)
	}

	logs.Push(*deployment)

	heap.Push(queue, *deployment)

	return int8(queue.Len())
}
