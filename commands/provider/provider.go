package provider

import (
	"github.com/mesg-foundation/core/daemon"
	"github.com/mesg-foundation/core/protobuf/coreapi"
)

// Provider is a struct that provides all methods required by any command.
type Provider struct {
	*CoreProvider
	*ServiceProvider
	*WorkflowProvider
}

// New creates Provider based on given CoreClient.
func New(c coreapi.CoreClient, d daemon.Daemon) *Provider {
	return &Provider{
		CoreProvider:     NewCoreProvider(c, d),
		ServiceProvider:  NewServiceProvider(c),
		WorkflowProvider: NewWorkflowProvider(c),
	}
}
