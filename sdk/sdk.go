package sdk

import (
	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	eventsdk "github.com/mesg-foundation/engine/sdk/event"
	executionsdk "github.com/mesg-foundation/engine/sdk/execution"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	ownershipsdk "github.com/mesg-foundation/engine/sdk/ownership"
	processesdk "github.com/mesg-foundation/engine/sdk/process"
	runnersdk "github.com/mesg-foundation/engine/sdk/runner"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
)

// SDK exposes all functionalities of MESG Engine.
type SDK struct {
	Service   *servicesdk.SDK
	Instance  *instancesdk.SDK
	Execution *executionsdk.SDK
	Event     *eventsdk.Event
	Process   *processesdk.SDK
	Ownership *ownershipsdk.SDK
	Runner    *runnersdk.SDK
}

// New creates a new SDK with given options.
func New(client *cosmos.Client, kb *cosmos.Keybase, container container.Container, engineName, port, ipfsEndpoint string) *SDK {
	ps := pubsub.New(0)
	serviceSDK := servicesdk.New(client)
	ownershipSDK := ownershipsdk.New(client)
	instanceSDK := instancesdk.New(client)
	runnerSDK := runnersdk.New(client, serviceSDK, instanceSDK, container, engineName, port, ipfsEndpoint)
	processSDK := processesdk.New(client)
	executionSDK := executionsdk.New(client, serviceSDK, instanceSDK, runnerSDK)
	eventSDK := eventsdk.New(ps, serviceSDK, instanceSDK)
	return &SDK{
		Service:   serviceSDK,
		Instance:  instanceSDK,
		Execution: executionSDK,
		Event:     eventSDK,
		Process:   processSDK,
		Ownership: ownershipSDK,
		Runner:    runnerSDK,
	}
}
