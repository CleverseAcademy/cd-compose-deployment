package config

const (
	PathGetDeploymentJTI        = "/deploy/nextJTI/:serviceName"
	PathAddDeployment           = "/deploy"
	PathGetLatestDeploymentInfo = "/deployment/:serviceName"
)

const (
	ErrorServiceNotFound string = "SERVICE_NOT_FOUND"
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
)
