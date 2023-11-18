package services

import (
	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/CleverseAcademy/cd-compose-deployment/usecases"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/pkg/errors"
)

type IArgsCreateDeployNewImageHandler struct {
	ServiceName string
	ComposeAPI  api.Service
}

func (s Service) SoyDeploy(args IArgsCreateDeployNewImageHandler) (*types.Project, error) {
	prj, err := s.ExecuteServiceDeployments.Execute(
		usecases.IArgsExecuteServiceDeployments{
			ComposeAPI:                args.ComposeAPI,
			ServiceName:               entities.ServiceName(args.ServiceName),
			LogDeploymentDoneEvent:    s.LogDeploymentDoneEvent,
			LogDeploymentFailureEvent: s.LogDeploymentFailureEvent,
			LogDeploymentSkippedEvent: s.LogDeploymentSkippedEvent,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "ExecuteServiceDeployments")
	}

	return prj, nil
}
