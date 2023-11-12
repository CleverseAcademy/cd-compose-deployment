package usecases

import (
	"fmt"

	"github.com/CleverseAcademy/cd-compose-deployment/entities"
	"github.com/CleverseAcademy/cd-compose-deployment/utils"
	"github.com/pkg/errors"

	"github.com/compose-spec/compose-go/types"
	"github.com/docker/compose/v2/pkg/api"
)

type composer struct {
	serviceName entities.ServiceName `json:"-"`
	deployment  *entities.Deployment `json:"-"`
	*types.Project
}

// createComposer initialize Composer object only after all validation, i.e. config checksum, is done.
// In order to prevent deployment with an obsoleted version of the configuration.
func createComposer(p *types.Project, s entities.ServiceName, d *entities.Deployment) (composer, error) {
	for idx, svc := range p.Services {
		if svc.Name == string(s) {
			omitImageCfg := svc
			omitImageCfg.Image = ""

			checksum, err := utils.Base64EncodedSha256(omitImageCfg)
			if err != nil {
				return composer{}, errors.Wrap(err, "Base64EncodedSha256 failed")
			}

			if checksum != d.CfgChecksum {
				return composer{}, fmt.Errorf("Project's checksum mismatch with the checksum value specified in a deployment")
			}

			c := composer{
				Project:     p,
				serviceName: s,
				deployment:  d,
			}

			p.Services[idx].Image = d.Image

			return c, nil
		}
	}

	return composer{}, fmt.Errorf("service %s not found", s)
}

func (c *composer) applyTo(composeService api.Service) error {
	err := composeService.Stop(c.deployment.GetCtx(), c.Name, api.StopOptions{
		Project:  c.Project,
		Services: []string{string(c.serviceName)},
	})
	if err != nil {
		return errors.Wrap(err, "compose stop failed")
	}

	err = composeService.Create(c.deployment.GetCtx(), c.Project, api.CreateOptions{
		Services:             []string{string(c.serviceName)},
		Recreate:             api.RecreateForce,
		RecreateDependencies: api.RecreateNever,
	})
	if err != nil {
		return errors.Wrap(err, "compose create failed")
	}

	err = composeService.Start(c.deployment.GetCtx(), c.Name, api.StartOptions{
		Project: c.Project,
	})

	return errors.Wrap(err, "compose start failed")
}
