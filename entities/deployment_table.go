package entities

import (
	"container/heap"
	"fmt"
	"reflect"
	"sync"

	"github.com/CleverseAcademy/cd-compose-deployment/constants"
)

type ServiceName string

type DeploymentQueue struct {
	deployments []Deployment
	order       bool
	sync.RWMutex
}
type DeploymentTable map[ServiceName]*DeploymentQueue

func NewDeploymentTable() *DeploymentTable {
	tbl := make(DeploymentTable)
	return &tbl
}

func (dq *DeploymentQueue) Len() int { return len(dq.deployments) }

func (dq *DeploymentQueue) Less(i, j int) bool {
	dq.RLock()
	defer dq.RUnlock()

	if dq.order == constants.DescOrder {
		return dq.deployments[i].Priority > dq.deployments[j].Priority
	}
	return dq.deployments[i].Priority < dq.deployments[j].Priority
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
	return append(make([]Deployment, 0), dq.deployments...)
}

func (dq *DeploymentQueue) At(i int) *Deployment {
	return &dq.deployments[i]
}

func (dq *DeploymentQueue) Copy(order bool) *DeploymentQueue {
	q := &DeploymentQueue{
		deployments: dq.Items(),
		order:       order,
	}

	heap.Init(q)

	return q
}

func (tbl DeploymentTable) GetServiceDeploymentQueue(serviceName ServiceName) (*DeploymentQueue, error) {
	fmt.Printf("%v\n", tbl)
	val, found := (tbl)[serviceName]
	if !found {
		return nil, fmt.Errorf("%s: %s", constants.ErrorEmptyDeployment, serviceName)
	}

	return val, nil
}

func (tbl DeploymentTable) InitializeDeploymentQueue(serviceName ServiceName) *DeploymentQueue {
	q := new(DeploymentQueue)
	q.deployments = make([]Deployment, 0)
	q.order = constants.DescOrder

	heap.Init(q)
	(tbl)[serviceName] = q

	return q
}
