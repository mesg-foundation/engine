package sdk

import (
	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/database"
	eventsdk "github.com/mesg-foundation/engine/sdk/event"
	executionsdk "github.com/mesg-foundation/engine/sdk/execution"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
)

// SDK exposes all functionalities of MESG Engine.
type SDK struct {
	Service   *servicesdk.Service
	Instance  *instancesdk.Instance
	Execution *executionsdk.Execution
	Event     *eventsdk.Event
}

// New creates a new SDK with given options.
func New(c container.Container, serviceDB database.ServiceDB, instanceDB database.InstanceDB, execDB database.ExecutionDB, engineName, port string) *SDK {
	ps := pubsub.New(0)
	serviceSDK := servicesdk.New(c, serviceDB)
	instanceSDK := instancesdk.New(c, serviceSDK, instanceDB, engineName, port)
	eventSDK := eventsdk.New(ps, serviceSDK, instanceSDK)
	executionSDK := executionsdk.New(ps, serviceSDK, instanceSDK, eventSDK, execDB)
	return &SDK{
		Service:   serviceSDK,
		Instance:  instanceSDK,
		Execution: executionSDK,
		Event:     eventSDK,
	}
}
