package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	ComposeWorkingDir  string
	ComposeFile        string
	ComposeProjectName string
	DockerContext      string
	PublicKeyPEMBytes  []byte
	InitialHash        string
}

var AppConfig Config

func init() {
	viper.SetEnvPrefix("CD")
	viper.AutomaticEnv()
	viper.SetDefault("COMPOSE_FILE", "/run/secrets/compose-file")
	viper.SetDefault("COMPOSE_PROJECT_NAME", "")
	viper.SetDefault("DOCKER_CONTEXT", "default")
	viper.SetDefault("KEYPAIR_PUBKEY_FILE", "./keypairs/ecpubkey.pem")
}

func init() {
	workingDir := viper.GetString("HOST_COMPOSE_WORKING_DIR")
	if len(workingDir) == 0 {
		panic("ENV: CD_HOST_COMPOSE_WORKING_DIR is not configured")
	}
	pem, err := os.ReadFile(viper.GetString("KEYPAIR_PUBKEY_FILE"))
	if err != nil {
		panic(err)
	}

	AppConfig = Config{
		ComposeWorkingDir:  workingDir,
		ComposeFile:        viper.GetString("COMPOSE_FILE"),
		ComposeProjectName: viper.GetString("COMPOSE_PROJECT_NAME"),
		DockerContext:      viper.GetString("DOCKER_CONTEXT"),
		PublicKeyPEMBytes:  pem,
		InitialHash:        "f8c0c5c0811c1344e6948c5fabc2839151cd7f0444c2724f2cddd238ce62bdec",
	}
}
