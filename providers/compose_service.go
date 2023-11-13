package providers

import (
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

func GetComposeService(clnt *client.Client, contextName string) (api.Service, error) {
	dockerCli, err := command.NewDockerCli(command.WithAPIClient(clnt))
	if err != nil {
		return nil, errors.Wrap(err, "create docker cli")
	}

	err = dockerCli.Initialize(&flags.ClientOptions{
		Context: contextName,
	})
	if err != nil {
		return nil, errors.Wrap(err, "initialize with context")
	}

	composeAPI := compose.NewComposeService(dockerCli)
	return composeAPI, nil
}
