package provider

import "github.com/mesg-foundation/core/interface/grpc/core"

type Provider struct {
	*CoreProvider
	*ServiceProvider
}

func New(c core.CoreClient) *Provider {
	return &Provider{
		CoreProvider:    NewCoreProvider(c),
		ServiceProvider: NewServiceProvider(c),
	}
}
