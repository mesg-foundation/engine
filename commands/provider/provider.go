package provider

import "github.com/mesg-foundation/core/interface/grpc/core"

// Provider is a struct that provides all methods required by any command.
type Provider struct {
	*ServiceProvider
}

// New creates Provider based on given CoreClient.
func New(c core.CoreClient) *Provider {
	return &Provider{
		ServiceProvider: NewServiceProvider(c),
	}
}
