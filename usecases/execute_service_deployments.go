package usecases

import (
	"container/heap"
	"context"
	"fmt"
	"reflect"

	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/compose/v2/pkg/api"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

type UseCaseExecuteServiceDeployments struct {
	*UseCaseEnqueueServiceDeployment
}

func (u *UseCaseExecuteServiceDeployments) Execute(clnt *client.Client, composeAPI api.Service, svcName entities.ServiceName) (*types.Project, error) {
	u.Lock()
	defer u.Unlock()

	queue, err := u.tbl.GetServiceDeploymentQueue(svcName)
	if err != nil {
		return nil, fmt.Errorf("ExecuteDeployment: %w", err)
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
		return nil, errors.Wrap(err, "UseCaseExecuteServiceDeployments@createComposer")
	}

	err = composer.applyTo(composeAPI)
	if err != nil {
		return nil, errors.Wrap(err, "UseCaseExecuteServiceDeployments@composer.applyTo")
	}

	logs, err := u.Logs.GetServiceDeploymentQueue(svcName)
	if err != nil {
		return nil, errors.Wrap(err, "UseCaseExecuteServiceDeployments@Logs.GetServiceDeploymentQueue")
	}

	composeLabels := filters.NewArgs(filters.KeyValuePair{
		Key:   "label",
		Value: fmt.Sprintf("%s=%s", api.ProjectLabel, u.Project.Name),
	}, filters.KeyValuePair{
		Key:   "label",
		Value: fmt.Sprintf("%s=%s", api.ServiceLabel, svcName),
	})

	containers, err := clnt.ContainerList(context.Background(), dockertypes.ContainerListOptions{
		Filters: composeLabels,
	})
	if err != nil {
		return nil, errors.Wrap(err, "UseCaseExecuteServiceDeployments@ContainerList")
	}

	for idx, v := range logs.Items() {
		if v.Ref == highestPDeployment.Ref && v.Timestamp.Equal(highestPDeployment.Timestamp) && v.Image == highestPDeployment.Image {
			logs.At(idx).Container = &containers[0]
		}
	}

	return &u.Project, nil
}
