package config

import "github.com/spf13/viper"

type Config struct {
	ComposeWorkingDir  string
	ComposeFile        string
	ComposeProjectName string
	DockerContext      string
}

var AppConfig Config

func init() {
	viper.SetEnvPrefix("CD")
	viper.AutomaticEnv()
	viper.SetDefault("COMPOSE_FILE", "/run/secrets/compose-file")
	viper.SetDefault("COMPOSE_PROJECT_NAME", "")
	viper.SetDefault("DOCKER_CONTEXT", "default")
}

func init() {
	workingDir := viper.GetString("HOST_COMPOSE_WORKING_DIR")
	if len(workingDir) == 0 {
		panic("ENV: CD_HOST_COMPOSE_WORKING_DIR is not configured")
	}
	AppConfig = Config{
		ComposeWorkingDir:  workingDir,
		ComposeFile:        viper.GetString("COMPOSE_FILE"),
		ComposeProjectName: viper.GetString("COMPOSE_PROJECT_NAME"),
		DockerContext:      viper.GetString("DOCKER_CONTEXT"),
	}
}
