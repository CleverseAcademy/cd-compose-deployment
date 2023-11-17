package usecases

import (
	"container/heap"
	"fmt"
	"reflect"

	"github.com/CleverseAcademy/cd-compose-deployment/constants"
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/pkg/errors"
)

type UseCaseExecuteServiceDeployments struct {
	*UseCaseEnqueueServiceDeployment
}

func (u *UseCaseExecuteServiceDeployments) Execute(composeAPI api.Service, svcName entities.ServiceName) (*types.Project, *entities.Deployment, error) {
	u.Lock()
	defer u.Unlock()

	if u.tbl == nil {
		return &u.Project, nil, fmt.Errorf("%s: %s", constants.ErrorEmptyDeployment, svcName)
	}

	queue, err := u.tbl.GetServiceDeploymentQueue(svcName)
	if err != nil {
		return &u.Project, nil, fmt.Errorf("ExecuteDeployment: %w", err)
	}
	if queue.Len() == 0 {
		return &u.Project, nil, fmt.Errorf("%s: %s", constants.ErrorEmptyDeployment, svcName)
	}

	highestPItem := heap.Pop(queue)
	highestPDeployment, ok := highestPItem.(entities.Deployment)
	if !ok {
		return &u.Project, nil, fmt.Errorf("Given deployment is of type %s, not entities.Deployment", reflect.TypeOf(highestPItem).String())
	}

	for queue.Len() > 0 {
		item := heap.Pop(queue)

		deployment, ok := item.(entities.Deployment)
		if !ok {
			panic(fmt.Sprintf("Given deployment is of type %s, not entities.Deployment", reflect.TypeOf(item).String()))
		}

		deployment.Cancel()
	}

	composer, err := createComposer(
		&u.Project,
		svcName,
		&highestPDeployment,
	)
	if err != nil {
		return &u.Project, &highestPDeployment, errors.Wrap(err, "UseCaseExecuteServiceDeployments@createComposer")
	}

	err = composer.applyTo(composeAPI)
	if err != nil {
		return &u.Project, &highestPDeployment, errors.Wrap(err, "UseCaseExecuteServiceDeployments@composer.applyTo")
	}

	u.tbl = nil
	return &u.Project, &highestPDeployment, nil
}
