package usecases

import (
	"io"
	"sync"

	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/docker/client"
)

type DeploymentUseCase struct {
	// Project is intended to copy orignal value to its own
	Project types.Project
	sync.RWMutex
}

type EventLogUseCase struct {
	Logger       io.Writer
	DockerClient *client.Client
}

type (
	IUseCasePrepareServiceDeployment interface {
		Execute(service entities.ServiceName, priotity int8, ref string, image string) (*entities.Deployment, error)
	}
	IUseCaseEnqueueServiceDeployment interface {
		Execute(service entities.ServiceName, deployment *entities.Deployment) int8
	}
	IUseCaseExecuteServiceDeployments interface {
		Execute(composeAPI api.Service, service entities.ServiceName) (*types.Project, *entities.Deployment, error)
	}
	IUseCaseGetAllServiceDeploymentInfo interface {
		Execute(service entities.ServiceName) ([]entities.Deployment, error)
	}
	IUseCaseGetServiceDeploymentInfo interface {
		Execute(service entities.ServiceName, ref string) (entities.Deployment, error)
	}
	IUseCaseGetLatestServiceDeploymentInfo interface {
		Execute(service entities.ServiceName) (entities.Deployment, error)
	}
)

type (
	IUseCaseLogConfigLoadedEvent interface {
		Execute(types.Project) error
	}
	IUseCaseLogConfigChangesDetectedEvent interface {
		Execute(types.Project) error
	}
	IUseCaseLogStopSignalReceivedEvent interface {
		Execute(types.Project) error
	}
	IUseCaseLogDeploymentDoneEvent interface {
		Execute(types.Project, entities.Deployment, entities.ServiceName) error
	}
	IUseCaseLogDeploymentFailureEvent interface {
		Execute(types.Project, entities.UndeployableServiceInfo) error
	}
	IUseCaseLogDeploymentSkippedEvent interface {
		Execute(entities.ServiceName) error
	}
)
