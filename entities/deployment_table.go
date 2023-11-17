package entities

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/CleverseAcademy/cd-compose-deployment/constants"
)

type ServiceName string

type DeploymentQueue struct {
	deployments []Deployment
	sync.RWMutex
}
type DeploymentTable map[ServiceName]*DeploymentQueue

func (dq *DeploymentQueue) Len() int { return len(dq.deployments) }

func (dq *DeploymentQueue) Less(i, j int) bool {
	dq.RLock()
	defer dq.RUnlock()

	return dq.deployments[i].Priority > dq.deployments[j].Priority
}

func (dq *DeploymentQueue) Swap(i, j int) {
	dq.Lock()
	defer dq.Unlock()

	dq.deployments[i], dq.deployments[j] = dq.deployments[j], dq.deployments[i]
}

func (dq *DeploymentQueue) Push(x any) {
	dq.Lock()
	defer dq.Unlock()

	asserted, ok := x.(Deployment)
	if !ok {
		panic(fmt.Sprintf("Given x is of type %s, not Deployment", reflect.TypeOf(x).String()))
	}

	dq.deployments = append(dq.deployments, asserted)
}

func (dq *DeploymentQueue) Pop() any {
	dq.Lock()
	defer dq.Unlock()

	old := dq.deployments
	n := len(old)
	deployment := old[n-1]
	dq.deployments = old[0 : n-1]

	return deployment
}

func (dq *DeploymentQueue) Items() []Deployment {
	return append([]Deployment{}, dq.deployments...)
}

func (dq *DeploymentQueue) At(i int) *Deployment {
	return &dq.deployments[i]
}

func (tbl DeploymentTable) GetServiceDeploymentQueue(serviceName ServiceName) (*DeploymentQueue, error) {
	val, found := (tbl)[serviceName]
	if !found {
		return nil, fmt.Errorf("%s DeploymentQueue for service %s not found", constants.ErrorEmptyDeployment, serviceName)
	}

	return val, nil
}

func (tbl DeploymentTable) InitializeDeploymentQueue(serviceName ServiceName) *DeploymentQueue {
	q := &DeploymentQueue{}
	(tbl)[serviceName] = q

	return q
}
