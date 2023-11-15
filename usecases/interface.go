package usecases

import (
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

type (
	IUseCasePrepareServiceDeployment interface {
		Execute(service entities.ServiceName, priotity int8, ref string, image string) (*entities.Deployment, error)
	}
	IUseCaseEnqueueServiceDeployment interface {
		Execute(service entities.ServiceName, deployment *entities.Deployment) int8
	}
	IUseCaseExecuteServiceDeployments interface {
		Execute(clnt *client.Client, composeAPI api.Service, service entities.ServiceName) (*types.Project, error)
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
