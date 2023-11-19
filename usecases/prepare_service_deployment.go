package usecases

import (
	"fmt"

	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/pkg/errors"
)

type UseCasePrepareServiceDeployment struct {
	*DeploymentUseCase
}

func (u *UseCasePrepareServiceDeployment) Execute(service entities.ServiceName, priotity int8, ref string, image string) (*entities.Deployment, error) {
	for idx, svc := range u.Project.Services {
		if svc.Name == string(service) {
			u.Lock()
			defer u.Unlock()

			target := u.Project.Services[idx]
			target.Image = image

			deployment, err := entities.CreateDeployment(priotity, ref, &target)
			if err != nil {
				return nil, errors.Wrap(err, "entities.CreateDeployment")
			}

			return deployment, nil
		}
	}

	return nil, fmt.Errorf("service %s not found in the project", string(service))
}
