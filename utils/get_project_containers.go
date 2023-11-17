package utils

import (
	"context"
	"fmt"

	"github.com/docker/compose/v2/pkg/api"
	docker "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func GetProjectContainers(prjName string, clnt *client.Client) ([]docker.Container, error) {
	composeLabels := filters.NewArgs(
		filters.KeyValuePair{
			Key:   "label",
			Value: fmt.Sprintf("%s=%s", api.ProjectLabel, prjName),
		},
	)

	containers, err := clnt.ContainerList(context.Background(), docker.ContainerListOptions{
		Filters: composeLabels,
		All:     true,
	})
	return containers, errors.Wrap(err, "GetProjectContainers")
}
