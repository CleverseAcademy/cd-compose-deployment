package services

import (
	"github.com/CleverseAcademy/cd-compose-deployment/usecases"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/docker/client"
)

type Service struct {
	ExecuteServiceDeployments usecases.IUseCaseExecuteServiceDeployments
	LogDeploymentDoneEvent    usecases.IUseCaseLogDeploymentDoneEvent
	LogDeploymentFailureEvent usecases.IUseCaseLogDeploymentFailureEvent
	LogDeploymentSkippedEvent usecases.IUseCaseLogDeploymentSkippedEvent
}

type IService interface {
	PeriodicallySoyDeploy(*client.Client, api.Service, *types.Project, string)
}
