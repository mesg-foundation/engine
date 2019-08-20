package servicesdk

import (
	"fmt"

	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
	"github.com/mesg-foundation/engine/cosmos"
)

type Cosmos struct {
}

func NewCosmos(app *cosmos.App, c container.Container, keeperFactory KeeperFactor) *Cosmos {
	app.RegisterModule(newAppModule(c, keeperFactory))
	return &Cosmos{}
}

// Create creates a new service from definition.
func (s *Cosmos) Create(srv *service.Service) (*service.Service, error) {
	return nil, fmt.Errorf("not implemented")
}

// Delete deletes the service by hash.
func (s *Cosmos) Delete(hash hash.Hash) error {
	return fmt.Errorf("not implemented")
}

// Get returns the service that matches given hash.
func (s *Cosmos) Get(hash hash.Hash) (*service.Service, error) {
	return nil, fmt.Errorf("not implemented")
}

// List returns all services.
func (s *Cosmos) List() ([]*service.Service, error) {
	return nil, fmt.Errorf("not implemented")
}
