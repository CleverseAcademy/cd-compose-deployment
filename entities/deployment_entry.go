package entities

import (
	"context"
	"time"

	"github.com/CleverseAcademy/cd-compose-deployment/utils"
	composetypes "github.com/compose-spec/compose-go/types"
	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
)

type StatusEntry struct {
	Timestamp   time.Time        `json:"ts"`
	Image       string           `json:"image"`
	CfgChecksum string           `json:"cfg_checksum"`
	Container   *types.Container `json:"container,omitempty"`
}

type Deployment struct {
	Priority  int8               `json:"p"`
	Ref       string             `json:"ref"`
	ctx       context.Context    `json:"-"`
	abortFunc context.CancelFunc `json:"-"`
	*StatusEntry
}

func CreateDeployment(p int8, ref string, cfg *composetypes.ServiceConfig) (*Deployment, error) {
	image := cfg.Image

	omitImageCfg := *cfg
	omitImageCfg.Image = ""

	checksum, err := utils.Base64EncodedSha256(omitImageCfg)
	if err != nil {
		return nil, errors.Wrap(err, "CreateDeployment@utils.Base64EncodedSha256")
	}

	newEntry := &StatusEntry{
		Image:       image,
		Timestamp:   time.Now(),
		CfgChecksum: checksum,
		Container:   nil,
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Deployment{
		Priority:    p,
		Ref:         ref,
		StatusEntry: newEntry,
		abortFunc:   cancel,
		ctx:         ctx,
	}, nil
}

func (d *Deployment) GetCtx() context.Context {
	return d.ctx
}

func (d *Deployment) Cancel() {
	d.Container = &types.Container{
		Status: "Aborted",
		State:  "Aborted",
	}

	d.abortFunc()
}

func (d *Deployment) ID() (string, error) {
	normalized := *d
	normalized.Container = nil

	checksum, err := utils.Base64EncodedSha256(normalized)
	if err != nil {
		return "", errors.Wrap(err, "Deployment.ID@utils.Base64EncodedSha256")
	}

	return checksum, nil
}

func (d *Deployment) Checksum() (string, error) {
	checksum, err := utils.Base64EncodedSha256(d)
	if err != nil {
		return "", errors.Wrap(err, "Deployment.Checksum@utils.Base64EncodedSha256")
	}

	return checksum, nil
}
