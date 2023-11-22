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
	params := IArgsLoadComposeProject{
		WorkingDir: strings.TrimRight(args.WorkingDir, constants.PathSeperator),
	}

	if len(args.ProjectName) > 0 {
		params.ProjectName = args.ProjectName
	} else {
		chunks := strings.Split(params.WorkingDir, constants.PathSeperator)
		params.ProjectName = chunks[len(chunks)-1]
	}

	params.ComposeFile = args.ComposeFile

	options := func(o *loader.Options) {
		o.SetProjectName(params.ProjectName, true)
		o.Interpolate.LookupValue = os.LookupEnv
	}

	prj, err := loader.Load(types.ConfigDetails{
		WorkingDir:  params.WorkingDir,
		ConfigFiles: types.ToConfigFiles([]string{params.ComposeFile}),
	}, options)
	if err != nil {
		return nil, errors.Wrap(err, "loader.Load")
	}

	for idx, v := range prj.Services {
		if v.Labels == nil {
			prj.Services[idx].Labels = make(types.Labels)
		}
		prj.Services[idx].Labels[api.ProjectLabel] = params.ProjectName
		prj.Services[idx].Labels[api.ServiceLabel] = v.Name
		prj.Services[idx].Labels[api.OneoffLabel] = "False"
		prj.Services[idx].Labels[constants.ComposeDeploymentLabel] = "True"
	}

	return prj, nil
}
