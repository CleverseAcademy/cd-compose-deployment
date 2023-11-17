package entities

import (
	"time"

	"github.com/docker/docker/api/types"
)

type EventLog struct {
	Event           string    `json:"e"`
	Timestamp       time.Time `json:"ts"`
	ProjectChecksum string    `json:"prj_chksm"`
}

type ServiceInfo struct {
	Name        ServiceName `json:"service"`
	Image       string      `json:"image"`
	ContainerID string      `json:"ctnr_id,omitempty"`
}

type UndeployableServiceInfo struct {
	Name          ServiceName `json:"name"`
	Err           string      `json:"err"`
	DeploymentRef string      `json:"d_ref"`
	CfgChecksum   string      `json:"cfg_chksm"`
	Image         string      `json:"image"`
}

type DetailedServiceInfo struct {
	*ServiceDetail `json:",inline"`

	Name          ServiceName      `json:"name"`
	DeploymentRef string           `json:"d_ref,omitempty"`
	Container     *types.Container `json:"container"`
}

type StopSignalReceivedEventEntry struct {
	EventLog `json:",inline"`
}

type ConfigLoadedEventEntry struct {
	EventLog `json:",inline"`

	Services []ServiceInfo `json:"services"`
}

type ConfigChangesDetectedEventEntry struct {
	EventLog `json:",inline"`
}

type DeploymentDoneEventEntry struct {
	EventLog `json:",inline"`

	Services []DetailedServiceInfo `json:"services"`
}

type DeploymentSkippedEventEntry struct {
	EventLog `json:",inline"`

	Name ServiceName `json:"name"`
}

type DeploymentFailureEventEntry struct {
	EventLog      `json:",inline"`
	FailedService UndeployableServiceInfo `json:"failed_service"`

	Services []DetailedServiceInfo `json:"services"`
}
