package services

import (
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

type IArgsCreateDeployNewImageHandler struct {
	ServiceName string
	ComposeAPI  api.Service
	DockerClnt  *client.Client
}

func (s Service) SoyDeploy(args IArgsCreateDeployNewImageHandler) (*types.Project, error) {
	prj, err := s.ExecuteServiceDeployments.Execute(args.DockerClnt, args.ComposeAPI, entities.ServiceName(args.ServiceName))
	if err != nil {
		return nil, errors.Wrap(err, "SoyDeploy")
	}
	return prj, nil
}
