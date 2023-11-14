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
	fmt.Printf("Applying %s\n", c.serviceName)
	err := composeService.Up(c.deployment.GetCtx(), c.Project, api.UpOptions{
		Create: api.CreateOptions{
			Recreate:             api.RecreateDiverged,
			RecreateDependencies: api.RecreateNever,
			Services:             []string{string(c.serviceName)},
		},
	})
	if err != nil {
		return errors.Wrap(err, "compose up failed")
	}

	return nil
}
