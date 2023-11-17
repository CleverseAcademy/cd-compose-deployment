package constants

const (
	PathGetDeploymentJTI        = "/deploy/nextJTI/:serviceName"
	PathAddDeployment           = "/deploy"
	PathGetLatestDeploymentInfo = "/deployment/:serviceName"
)

const (
	ErrorServiceNotFound string = "SERVICE_NOT_FOUND"
	ErrorEmptyDeployment string = "EMPTY_DEPLOYMENT"
)

const (
	DeploymentDoneEventName        = "DEPLOYMENT_DONE"
	DeploymentFailureEventName     = "DEPLOYMENT_FAILURE"
	DeploymentSkippedEventName     = "DEPLOYMENT_SKIPPED"
	ConfigLoadedEventName          = "CONFIG_LOADED"
	ConfigChangesDetectedEventName = "CONFIG_CHANGES_DETECTED"
	StopSignalReceivedEventName    = "STOP_SIGNAL_RECEIVED"
)

const (
	EventLogDir = "soy_events"
)
