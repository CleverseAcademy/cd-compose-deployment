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

type IArgsExecuteServiceDeployments struct {
	ComposeAPI                api.Service
	ServiceName               entities.ServiceName
	LogDeploymentDoneEvent    IUseCaseLogDeploymentDoneEvent
	LogDeploymentFailureEvent IUseCaseLogDeploymentFailureEvent
	LogDeploymentSkippedEvent IUseCaseLogDeploymentSkippedEvent
}

func (u *UseCaseExecuteServiceDeployments) Execute(args IArgsExecuteServiceDeployments) (*types.Project, error) {
	u.Lock()
	defer u.Unlock()

	logSkipped := func() {
		loggingErr := args.LogDeploymentSkippedEvent.Execute(args.ServiceName)
		if loggingErr != nil {
			panic(loggingErr)
		}
	}

	logFailure := func(d *entities.Deployment, err error) {
		info := entities.UndeployableServiceInfo{
			Name:          args.ServiceName,
			DeploymentRef: "",
			Err:           err.Error(),
			CfgChecksum:   "",
			Image:         "",
		}

		if d != nil {
			info.DeploymentRef = d.Ref
			info.CfgChecksum = d.CfgChecksum
			info.Image = d.Image
		}

		loggingErr := args.LogDeploymentFailureEvent.Execute(u.Project, info)
		if loggingErr != nil {
			panic(loggingErr)
		}
		fmt.Printf("after usecases project: %v", u.Project)
	}

	if u.tbl == nil {
		defer logSkipped()

		return &u.Project, fmt.Errorf("%s: %s", constants.ErrorEmptyDeployment, args.ServiceName)
	}

	queue, err := u.tbl.GetServiceDeploymentQueue(args.ServiceName)
	if err != nil {
		defer logSkipped()

		return &u.Project, fmt.Errorf("GetServiceDeploymentQueue: %w", err)
	}
	if queue.Len() == 0 {
		defer logSkipped()

		return &u.Project, fmt.Errorf("%s: %s", constants.ErrorEmptyDeployment, args.ServiceName)
	}

	highestPItem := heap.Pop(queue)
	highestPDeployment, ok := highestPItem.(entities.Deployment)
	if !ok {
		panic(fmt.Errorf("Given deployment is of type %s, not entities.Deployment", reflect.TypeOf(highestPItem).String()))
	}

	for queue.Len() > 0 {
		item := heap.Pop(queue)

		deployment, ok := item.(entities.Deployment)
		if !ok {
			panic(fmt.Sprintf("Given deployment is of type %s, not entities.Deployment", reflect.TypeOf(item).String()))
		}

		deployment.Cancel()
	}

	fmt.Printf("before usecases project: %v", u.Project)
	composer, err := createComposer(
		u.Project,
		args.ServiceName,
		&highestPDeployment,
	)
	if err != nil {
		wrappedErr := errors.Wrap(err, "createComposer")
		defer logFailure(&highestPDeployment, wrappedErr)

		return &u.Project, wrappedErr
	}

	err = composer.applyTo(args.ComposeAPI)
	if err != nil {
		wrappedErr := errors.Wrap(err, "composer.applyTo")
		defer logFailure(&highestPDeployment, wrappedErr)

		return &u.Project, wrappedErr
	}

	u.tbl = nil

	err = args.LogDeploymentDoneEvent.Execute(composer.Project, highestPDeployment, args.ServiceName)
	if err != nil {
		panic(err)
	}

	return &composer.Project, nil
}
