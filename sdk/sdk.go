package sdk

import (
	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/database"
	eventsdk "github.com/mesg-foundation/engine/sdk/event"
	executionsdk "github.com/mesg-foundation/engine/sdk/execution"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	workflowsdk "github.com/mesg-foundation/engine/sdk/workflow"
)

// SDK exposes all functionalities of MESG Engine.
type SDK struct {
	Service   servicesdk.Service
	Instance  *instancesdk.Instance
	Execution *executionsdk.Execution
	Event     *eventsdk.Event
	Workflow  *workflowsdk.Workflow
}

// New creates a new SDK with given options.
func New(c container.Container, serviceSDK servicesdk.Service, instanceDB database.InstanceDB, execDB database.ExecutionDB, workflowDB database.WorkflowDB, engineName, port string) *SDK {
	ps := pubsub.New(0)
	instanceSDK := instancesdk.New(c, serviceSDK, instanceDB, engineName, port)
	workflowSDK := workflowsdk.New(instanceSDK, workflowDB)
	executionSDK := executionsdk.New(ps, serviceSDK, instanceSDK, workflowSDK, execDB)
	eventSDK := eventsdk.New(ps, serviceSDK, instanceSDK)
	return &SDK{
		Service:   serviceSDK,
		Instance:  instanceSDK,
		Execution: executionSDK,
		Event:     eventSDK,
		Workflow:  workflowSDK,
	}
}
