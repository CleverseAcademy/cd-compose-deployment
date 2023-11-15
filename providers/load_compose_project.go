package providers

import (
	"fmt"
	"os"
	"strings"

	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/pkg/errors"
)

const pathSeperator = "/"

type IArgsLoadComposeProject struct {
	WorkingDir  string
	ComposeFile string
	ProjectName string
}

func LoadComposeProject(args IArgsLoadComposeProject) (*types.Project, error) {
	params := IArgsLoadComposeProject{
		WorkingDir: strings.TrimRight(args.WorkingDir, pathSeperator),
	}

	chunks := strings.Split(params.WorkingDir, pathSeperator)

	if len(args.ProjectName) > 0 {
		params.ProjectName = args.ProjectName
	} else {
		params.ProjectName = chunks[len(chunks)-1]
	}

	if len(args.ComposeFile) > 0 {
		params.ComposeFile = args.ComposeFile
	} else {
		params.ComposeFile = fmt.Sprintf("%s/compose.yml", params.WorkingDir)
	}

	options := func(o *loader.Options) {
		o.SetProjectName(params.ProjectName, true)
		o.Interpolate.LookupValue = os.LookupEnv
	}

	prj, err := loader.Load(types.ConfigDetails{
		WorkingDir:  params.WorkingDir,
		ConfigFiles: types.ToConfigFiles([]string{params.ComposeFile}),
	}, options)
	if err != nil {
		return nil, errors.Wrap(err, "LoadComposeProject@loader.Load")
	}

	for idx, v := range prj.Services {
		if v.Labels == nil {
			prj.Services[idx].Labels = make(types.Labels)
		}
		prj.Services[idx].Labels[api.ProjectLabel] = params.ProjectName
		prj.Services[idx].Labels[api.ServiceLabel] = v.Name
		prj.Services[idx].Labels[api.OneoffLabel] = "False"
	}
	return prj, nil
}
