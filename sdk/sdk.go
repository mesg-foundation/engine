package sdk

import (
	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	eventsdk "github.com/mesg-foundation/engine/sdk/event"
	executionsdk "github.com/mesg-foundation/engine/sdk/execution"
	processesdk "github.com/mesg-foundation/engine/sdk/process"
	runnersdk "github.com/mesg-foundation/engine/sdk/runner"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
)

// SDK exposes all functionalities of MESG Engine.
type SDK struct {
	Service   *servicesdk.SDK
	Execution *executionsdk.SDK
	Event     *eventsdk.Event
	Process   *processesdk.SDK
	Runner    *runnersdk.SDK
}

// New creates a new SDK with given options.
func New(client *cosmos.Client, kb *cosmos.Keybase, container container.Container, engineName, port, ipfsEndpoint string) *SDK {
	ps := pubsub.New(0)
	serviceSDK := servicesdk.New(client)
	runnerSDK := runnersdk.New(client, serviceSDK, container, engineName, port, ipfsEndpoint)
	processSDK := processesdk.New(client)
	executionSDK := executionsdk.New(client, serviceSDK, runnerSDK)
	eventSDK := eventsdk.New(ps, serviceSDK, client)
	return &SDK{
		Service:   serviceSDK,
		Execution: executionSDK,
		Event:     eventSDK,
		Process:   processSDK,
		Runner:    runnerSDK,
	}
}
