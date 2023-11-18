package services

import (
	"github.com/CleverseAcademy/cd-compose-deployment/usecases"
	"github.com/compose-spec/compose-go/types"
)

type Service struct {
	GetAllServiceDeploymentInfo usecases.IUseCaseGetAllServiceDeploymentInfo
	ExecuteServiceDeployments   usecases.IUseCaseExecuteServiceDeployments
	LogDeploymentDoneEvent      usecases.IUseCaseLogDeploymentDoneEvent
	LogDeploymentFailureEvent   usecases.IUseCaseLogDeploymentFailureEvent
	LogDeploymentSkippedEvent   usecases.IUseCaseLogDeploymentSkippedEvent
}

type IService interface {
	GetNextJTI(serviceName string) (string, error)
	SoyDeploy(args IArgsCreateDeployNewImageHandler) (*types.Project, error)
}
