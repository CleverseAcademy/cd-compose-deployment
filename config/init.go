package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

const (
	envComposeFile           = "COMPOSE_FILE"
	envComposeProjectName    = "COMPOSE_PROJECT_NAME"
	envHostComposeWorkingDir = "HOST_COMPOSE_WORKING_DIR"
	envDockerContext         = "DOCKER_CONTEXT"
	envPubkeyFile            = "PUBKEY_FILE"
	envInitialHash           = "INITIAL_HASH"
	envTokenWindow           = "TOKEN_WINDOW"
	envPortBinding           = "PORT_BINDING"
	envDeployInterval        = "DEPLOY_INTERVAL_SECONDS"
	envDatadir               = "DATA_DIR"
)

type Config struct {
	ListeningSocket    string
	ComposeWorkingDir  string
	ComposeFile        string
	ComposeProjectName string
	DockerContext      string
	PublicKeyPEMBytes  []byte
	InitialHash        string
	DataDir            string
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
	viper.SetDefault(envDeployInterval, 20)
	viper.SetDefault(envDatadir, "./data/")
}

func init() {
	workingDir := viper.GetString(envHostComposeWorkingDir)
	if len(workingDir) == 0 {
		panic("ENV: CD_" + envHostComposeWorkingDir + " is not configured")
	}

	// 6071  openssl genpkey -algorithm EC -out eckey.pem \
	//  -pkeyopt ec_paramgen_curve:P-256 \
	//  -pkeyopt ec_param_enc:named_curve
	// 6072  openssl pkey -in eckey.pem -pubout -out ecpubkey.pem
	keyAbsolutePath, err := filepath.Abs(viper.GetString(envPubkeyFile))
	if err != nil {
		panic(err)
	}

	pem, err := os.ReadFile(keyAbsolutePath)
	if err != nil {
		panic(err)
	}

	dataAbsolutePath, err := filepath.Abs(viper.GetString(envDatadir))
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
		DataDir:            dataAbsolutePath,
		TokenWindow:        time.Duration(viper.GetUint64(envTokenWindow)) * time.Second,
		DeployInterval:     time.Duration(viper.GetUint64(envDeployInterval)) * time.Second,
	}
}
