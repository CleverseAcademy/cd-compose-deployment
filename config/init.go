package config

import (
	"os"
	"time"

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
	TokenWindow        time.Duration
	DeployInterval     time.Duration
}

var AppConfig Config

func init() {
	viper.SetEnvPrefix("CD")
	viper.AutomaticEnv()
	viper.SetDefault(envComposeFile, "/run/secrets/compose-file")
	viper.SetDefault(envComposeProjectName, "")
	viper.SetDefault(envDockerContext, "default")
	viper.SetDefault(envPubkeyFile, "./keypairs/ecpubkey.pem")
	viper.SetDefault(envInitialHash, "f8c0c5c0811c1344e6948c5fabc2839151cd7f0444c2724f2cddd238ce62bdec")
	viper.SetDefault(envTokenWindow, 60)
	viper.SetDefault(envPortBinding, ":3000")
	viper.SetDefault(envDeployInterval, 15)
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
		ListeningSocket:    viper.GetString(envPortBinding),
		ComposeFile:        viper.GetString(envComposeFile),
		ComposeProjectName: viper.GetString(envComposeProjectName),
		DockerContext:      viper.GetString(envDockerContext),
		InitialHash:        viper.GetString(envInitialHash),
		TokenWindow:        time.Duration(viper.GetUint64(envTokenWindow)) * time.Second,
		DeployInterval:     time.Duration(viper.GetUint64(envDeployInterval)) * time.Second,
	}
}
