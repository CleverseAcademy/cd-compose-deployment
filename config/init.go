package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	ListeningSocket    string
	ComposeWorkingDir  string
	ComposeFile        string
	ComposeProjectName string
	DockerContext      string
	PublicKeyPEMBytes  []byte
	InitialHash        string
	TokenWindow        uint16
}

var AppConfig Config

func init() {
	viper.SetEnvPrefix("CD")
	viper.AutomaticEnv()
	viper.SetDefault("COMPOSE_FILE", "/run/secrets/compose-file")
	viper.SetDefault("COMPOSE_PROJECT_NAME", "")
	viper.SetDefault("DOCKER_CONTEXT", "default")
	viper.SetDefault("PUBKEY_FILE", "./keypairs/ecpubkey.pem")
	viper.SetDefault("INITIAL_HASH", "f8c0c5c0811c1344e6948c5fabc2839151cd7f0444c2724f2cddd238ce62bdec")
	viper.SetDefault("TOKEN_WINDOW", 60)
	viper.SetDefault("BINDING", ":3000")
}

func init() {
	workingDir := viper.GetString("HOST_COMPOSE_WORKING_DIR")
	if len(workingDir) == 0 {
		panic("ENV: CD_HOST_COMPOSE_WORKING_DIR is not configured")
	}

	// 6071  openssl genpkey -algorithm EC -out eckey.pem \
	//  -pkeyopt ec_paramgen_curve:P-256 \
	//  -pkeyopt ec_param_enc:named_curve
	// 6072  openssl pkey -in eckey.pem -pubout -out ecpubkey.pem
	pem, err := os.ReadFile(viper.GetString("PUBKEY_FILE"))
	if err != nil {
		panic(err)
	}

	AppConfig = Config{
		ComposeWorkingDir:  workingDir,
		PublicKeyPEMBytes:  pem,
		ListeningSocket:    viper.GetString("BINDING"),
		ComposeFile:        viper.GetString("COMPOSE_FILE"),
		ComposeProjectName: viper.GetString("COMPOSE_PROJECT_NAME"),
		DockerContext:      viper.GetString("DOCKER_CONTEXT"),
		InitialHash:        viper.GetString("INITIAL_HASH"),
		TokenWindow:        viper.GetUint16("TOKEN_WINDOW"),
	}
}
