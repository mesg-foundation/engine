package api

import (
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/systemservices"
)

// API exposes all functionalities of MESG core.
type API struct {
	db             database.ServiceDB
	execDB         database.ExecutionDB
	systemservices *systemservices.SystemServices
	container      container.Container
}

// Option is a configuration func for MESG.
type Option func(*API)

// New creates a new API with given options.
func New(db database.ServiceDB, execDB database.ExecutionDB, systemservices *systemservices.SystemServices, options ...Option) (*API, error) {
	a := &API{db: db, execDB: execDB, systemservices: systemservices}
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
