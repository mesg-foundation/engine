package api

import (
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/database"
)

// API exposes all functionalies of MESG core.
type API struct {
	db        *database.ServiceDB
	container *container.Container
	cfg       *config.Config
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
	a.cfg, err = config.Global()
	if err != nil {
		return nil, err
	}
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

// DatabaseOption returns an option to set a db.
func DatabaseOption(db *database.ServiceDB) Option {
	return func(a *API) {
		a.db = db
	}
}
