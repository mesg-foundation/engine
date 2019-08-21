package servicesdk

import (
	"context"
	"fmt"

	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
)

type sdk struct {
	container     container.Container
	keeperFactory KeeperFactory
}

func New(app *cosmos.App, c container.Container, keeperFactory KeeperFactory) Service {
	// app.RegisterModule(newAppModule(c, keeperFactory))
	return &sdk{
		container:     c,
		keeperFactory: keeperFactory,
	}
}

// Create creates a new service from definition.
func (s *sdk) Create(srv *service.Service) (*service.Service, error) {
	return nil, fmt.Errorf("not implemented")
}

// Delete deletes the service by hash.
func (s *sdk) Delete(hash hash.Hash) error {
	return fmt.Errorf("not implemented")
}

// Get returns the service that matches given hash.
func (s *sdk) Get(hash hash.Hash) (*service.Service, error) {
	return nil, fmt.Errorf("not implemented")
}

// List returns all services.
func (s *sdk) List() ([]*service.Service, error) {
	return nil, fmt.Errorf("not implemented")
}

// ------------------------------------
//   Handlers
// ------------------------------------

// Create creates a new service from definition.
func (s *sdk) create(ctx context.Context, srv *service.Service) (*service.Service, error) {
	return create(ctx, s.container, s.keeperFactory, srv)
}

// Delete deletes the service by hash.
func (s *sdk) delete(ctx context.Context, hash hash.Hash) error {
	keeper, err := s.keeperFactory(ctx)
	if err != nil {
		return err
	}
	return keeper.Delete(hash)
}

// Get returns the service that matches given hash.
func (s *sdk) get(ctx context.Context, hash hash.Hash) (*service.Service, error) {
	keeper, err := s.keeperFactory(ctx)
	if err != nil {
		return nil, err
	}
	return keeper.Get(hash)
}

// List returns all services.
func (s *sdk) list(ctx context.Context) ([]*service.Service, error) {
	keeper, err := s.keeperFactory(ctx)
	if err != nil {
		return nil, err
	}
	return keeper.All()
}
