package services

import "github.com/CleverseAcademy/cd-compose-deployment/usecases"

type Service struct {
	GetAllServiceDeploymentInfo usecases.IUseCaseGetAllServiceDeploymentInfo
}

type IService interface {
	GetNextJTI(serviceName string) (string, error)
}
