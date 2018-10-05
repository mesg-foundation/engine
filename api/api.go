package api

import (
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/database"
)

// API exposes all functionalities of MESG core.
type API struct {
	db        *database.ServiceDB
	container container.Container
}

// Option is a configuration func for MESG.
type Option func(*API)

// New creates a new API with given options.
func New(db *database.ServiceDB, options ...Option) (*API, error) {
	a := &API{db: db}
	for _, option := range options {
		option(a)
	}
	if a.container == nil {
		var err error
		a.container, err = container.New()
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

// ContainerOption configures underlying container access API.
func ContainerOption(container container.Container) Option {
	return func(a *API) {
		a.container = container
	}
}
