package constants

const DefaultComposeYMLFilename = "compose.yml"

const PathSeperator = "/"

const ComposeDeploymentLabel = "com.pose.soy.deployment"

const (
	PathGetJTI                = "/nextJTI"
	PathAddDeployment         = "/deploy"
	PathGetNextDeploymentInfo = "/deploy/latest/:serviceName"
)

const (
	ErrorServiceNotFound string = "SERVICE_NOT_FOUND"
	ErrorEmptyDeployment string = "EMPTY_DEPLOYMENT"
)

const (
	DeploymentDoneEventName        = "DPLY_DONE"
	DeploymentFailureEventName     = "DPLY_FAILURE"
	DeploymentSkippedEventName     = "DPLY_SKIPPED"
	ConfigLoadedEventName          = "CONFIG_LOADED"
	ConfigChangesDetectedEventName = "CONFIG_CHANGES_DETECTED"
	StopSignalReceivedEventName    = "STOP_SIGNAL_RECEIVED"
)

const (
	EventLogDir  = "soy_events"
	AccessLogDir = "soy_access"
)

const (
	AscOrder  = true
	DescOrder = false
)
