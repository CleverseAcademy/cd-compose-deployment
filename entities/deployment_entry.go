package entities

import (
	"context"
	"time"

	"github.com/CleverseAcademy/cd-compose-deployment/utils"
	composetypes "github.com/compose-spec/compose-go/types"
	"github.com/pkg/errors"
)

type ServiceDetail struct {
	Timestamp   time.Time `json:"ts"`
	Image       string    `json:"image"`
	CfgChecksum string    `json:"cfg_chksm"`
}

type Deployment struct {
	Priority  int8               `json:"p"`
	Ref       string             `json:"ref"`
	ctx       context.Context    `json:"-"`
	abortFunc context.CancelFunc `json:"-"`
	*ServiceDetail
}

func CreateDeployment(p int8, ref string, cfg *composetypes.ServiceConfig) (*Deployment, error) {
	image := cfg.Image

	omitImageCfg := *cfg
	omitImageCfg.Image = ""

	checksum, err := utils.Base64EncodedSha256(omitImageCfg)
	if err != nil {
		return nil, errors.Wrap(err, "utils.Base64EncodedSha256")
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Deployment{
		Priority: p,
		Ref:      ref,
		ServiceDetail: &ServiceDetail{
			Image:       image,
			Timestamp:   time.Now(),
			CfgChecksum: checksum,
		},
		abortFunc: cancel,
		ctx:       ctx,
	}, nil
}

func (d *Deployment) GetCtx() context.Context {
	return d.ctx
}

func (d *Deployment) Cancel() {
	d.abortFunc()
}
