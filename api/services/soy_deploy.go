package services

import (
	"strings"

	"github.com/CleverseAcademy/cd-compose-deployment/constants"
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
	prj, deployment, err := s.ExecuteServiceDeployments.Execute(args.DockerClnt, args.ComposeAPI, entities.ServiceName(args.ServiceName))
	if err != nil && strings.Contains(err.Error(), constants.ErrorEmptyDeployment) {
		loggingErr := s.LogDeploymentSkippedEvent.Execute(*prj, entities.ServiceName(args.ServiceName))
		if loggingErr != nil {
			panic(loggingErr)
		}

		return nil, errors.Wrap(err, "SoyDeploy")
	} else if err != nil {
		info := entities.UndeployableServiceInfo{
			Name:          entities.ServiceName(args.ServiceName),
			DeploymentRef: "",
			Err:           err.Error(),
			CfgChecksum:   "",
			Image:         "",
		}
		if deployment != nil {
			info.DeploymentRef = deployment.Ref
			info.CfgChecksum = deployment.CfgChecksum
			info.Image = deployment.Image
		}

		loggingErr := s.LogDeploymentFailureEvent.Execute(*prj, info)
		if loggingErr != nil {
			panic(loggingErr)
		}

		return nil, errors.Wrap(err, "SoyDeploy")
	}

	err = s.LogDeploymentDoneEvent.Execute(*prj, *deployment, entities.ServiceName(args.ServiceName))
	if err != nil {
		return nil, errors.Wrap(err, "SoyDeploy@LogDeploymentDoneEvent")
	}

	return prj, nil
}
