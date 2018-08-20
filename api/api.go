package api

import "github.com/mesg-foundation/core/container"

// API exposes all functionalies of MESG core.
type API struct {
	container *container.Container
}

// Option is a configuration func for MESG.
type Option func(*API)

// New creates a new API with given options.
func New(options ...Option) (*API, error) {
	a := &API{}
	for _, option := range options {
		option(a)
	}
	var err error
	if a.container == nil {
		a.container, err = container.New()
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

// ContainerOption configures underlying container access API.
func ContainerOption(container *container.Container) Option {
	return func(a *API) {
		a.container = container
	}
}
