package sdk

import (
	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/database"
	eventsdk "github.com/mesg-foundation/core/sdk/event"
	executionsdk "github.com/mesg-foundation/core/sdk/execution"
	instancesdk "github.com/mesg-foundation/core/sdk/instance"
	servicesdk "github.com/mesg-foundation/core/sdk/service"
)

// SDK exposes all functionalities of MESG core.
type SDK struct {
	Service   *servicesdk.Service
	Instance  *instancesdk.Instance
	Execution *executionsdk.Execution
	Event     *eventsdk.Event
}

// New creates a new SDK with given options.
func New(c container.Container, serviceDB database.ServiceDB, instanceDB database.InstanceDB, execDB database.ExecutionDB) *SDK {
	ps := pubsub.New(0)
	serviceSDK := servicesdk.New(c, serviceDB)
	instanceSDK := instancesdk.New(c, serviceSDK, instanceDB)
	executionSDK := executionsdk.New(ps, serviceSDK, instanceSDK, execDB)
	eventSDK := eventsdk.New(ps, serviceSDK)
	return &SDK{
		Service:   serviceSDK,
		Instance:  instanceSDK,
		Execution: executionSDK,
		Event:     eventSDK,
	}
}
