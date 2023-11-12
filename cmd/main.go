package main

import (
	"context"
	"os"

	"github.com/compose-spec/compose-go/loader"
	composetypes "github.com/compose-spec/compose-go/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/docker/docker/client"
)

func main() {
	clnt, err := client.NewClientWithOpts(client.WithHostFromEnv(), client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
	// 	All: true,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	dockerCli, err := command.NewDockerCli(command.WithAPIClient(clnt))
	if err != nil {
		panic(err)
	}

	err = dockerCli.Initialize(&flags.ClientOptions{
		Context: "desktop-linux",
	})

	if err != nil {
		panic(err)
	}

	workingDir := "/Users/intaniger/works/focusing/CleverseAcademy/learnhub-api"
	configFile := "/Users/intaniger/works/focusing/CleverseAcademy/learnhub-api/compose.yaml"

	composeAPI := compose.NewComposeService(dockerCli)
	prj, err := loader.Load(composetypes.ConfigDetails{
		WorkingDir:  workingDir,
		ConfigFiles: composetypes.ToConfigFiles([]string{configFile}),
	}, func(o *loader.Options) {
		o.SetProjectName("learnhub-api", true)
		o.Interpolate.LookupValue = os.LookupEnv
	})
	if err != nil {
		panic(err)
	}

	// err = composeAPI.Down(context.Background(), "learnhub-api", api.DownOptions{
	// 	RemoveOrphans: true,
	// 	Project:       prj,
	// 	Services:      []string{"server", "db"},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	for idx, v := range prj.Services {
		if v.Labels == nil {
			prj.Services[idx].Labels = make(composetypes.Labels)
		}
		prj.Services[idx].Labels[api.ProjectLabel] = prj.Name
		prj.Services[idx].Labels[api.ServiceLabel] = v.Name
		prj.Services[idx].Labels[api.OneoffLabel] = "False"
	}

	// err = composeAPI.Build(context.Background(), prj, api.BuildOptions{
	// 	Pull:     true,
	// 	Services: []string{"server", "db"},
	// })
	// if err != nil {
	// 	panic(err)
	// }

	err = composeAPI.Up(context.Background(), prj, api.UpOptions{
		Create: api.CreateOptions{
			Services: []string{"server", "db"},
			Recreate: api.RecreateForce,
		},
	})
	if err != nil {
		panic(err)
	}

	// filter := filters.NewArgs(filters.KeyValuePair{
	// 	Key:   "name",
	// 	Value: prj.Name,
	// })
	// containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
	// 	All:     true,
	// 	Filters: filter,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	// err = composeAPI.Start(context.Background(), prj.Name, api.StartOptions{
	// 	Project:  prj,
	// 	Services: []string{"server", "db"},
	// })
	// if err != nil {
	// 	fmt.Printf("%v", err)
	// }

	// Start: api.StartOptions{
	// 	Project:  prj,
	// 	Services: []string{"server", "db"},
	// },

	// app := fiber.New()

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hi there")
	// })

	// app.Listen(":3000")
}
