package providers

import (
	"os"
	"strings"

	"github.com/CleverseAcademy/cd-compose-deployment/constants"
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/pkg/errors"
)

type IArgsLoadComposeProject struct {
	WorkingDir  string
	ComposeFile string
	ProjectName string
}

func LoadComposeProject(args IArgsLoadComposeProject) (*types.Project, error) {
	options := func(o *loader.Options) {
		o.SetProjectName(args.ProjectName, true)
		o.Interpolate.LookupValue = os.LookupEnv
	}

	prj, err := loader.Load(types.ConfigDetails{
		WorkingDir:  strings.TrimRight(args.WorkingDir, constants.PathSeperator),
		ConfigFiles: types.ToConfigFiles([]string{args.ComposeFile}),
	}, options)
	if err != nil {
		return nil, errors.Wrap(err, "loader.Load")
	}

	for idx, v := range prj.Services {
		if v.Labels == nil {
			prj.Services[idx].Labels = make(types.Labels)
		}
		prj.Services[idx].Labels[api.ProjectLabel] = args.ProjectName
		prj.Services[idx].Labels[api.ServiceLabel] = v.Name
		prj.Services[idx].Labels[api.OneoffLabel] = "False"
		prj.Services[idx].Labels[constants.ComposeDeploymentLabel] = "True"
	}

	return prj, nil
}
