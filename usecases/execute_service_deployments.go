package usecases

import (
	"container/heap"
	"fmt"
	"reflect"

	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/pkg/errors"
)

type UseCaseExecuteServiceDeployments struct {
	*UseCaseEnqueueServiceDeployment
}

func (u *UseCaseExecuteServiceDeployments) Execute(composeAPI api.Service, svcName entities.ServiceName) (*types.Project, error) {
	queue, err := u.tbl.GetServiceDeploymentQueue(svcName)
	if err != nil {
		return nil, fmt.Errorf("Service %s not found", svcName)
	}

	highestPItem := heap.Pop(queue)
	highestPDeployment, ok := highestPItem.(entities.Deployment)
	if !ok {
		return nil, fmt.Errorf("Given deployment is of type %s, not entities.Deployment", reflect.TypeOf(highestPItem).String())
	}

	for queue.Len() > 0 {
		item := heap.Pop(queue)

		deployment, ok := item.(entities.Deployment)
		if !ok {
			panic(fmt.Sprintf("Given deployment is of type %s, not entities.Deployment", reflect.TypeOf(item).String()))
		}

		if deployment.Container == nil {
			deployment.Cancel()
		}
	}

	composer, err := createComposer(
		&u.Project,
		svcName,
		&highestPDeployment,
	)
	if err != nil {
		return nil, errors.Wrap(err, "CreateComposer failed")
	}

	err = composer.applyTo(composeAPI)
	if err != nil {
		return nil, err
	}

	return &u.Project, nil
}
