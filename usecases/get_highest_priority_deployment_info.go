package usecases

import (
	"container/heap"
	"fmt"
	"reflect"

	"github.com/CleverseAcademy/cd-compose-deployment/constants"
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/pkg/errors"
)

type UseCaseGetCurrentHighestPriorityDeploymentInfo struct {
	*UseCaseEnqueueServiceDeployment
}

func (u *UseCaseGetCurrentHighestPriorityDeploymentInfo) Execute(service entities.ServiceName) (*entities.Deployment, error) {
	u.RLock()
	defer u.RUnlock()

	queue, err := u.tbl.GetServiceDeploymentQueue(service)
	if err != nil {
		return nil, errors.Wrap(err, "GetServiceDeploymentQueue")
	}
	if queue.Len() == 0 {
		return nil, fmt.Errorf("%s: %s", constants.ErrorEmptyDeployment, service)
	}

	highestPItem := heap.Pop(queue)
	highestPDeployment, ok := highestPItem.(entities.Deployment)
	if !ok {
		panic(fmt.Errorf("Given deployment is of type %s, not entities.Deployment", reflect.TypeOf(highestPItem).String()))
	}
	defer heap.Push(queue, highestPDeployment)

	return &highestPDeployment, nil
}
