package servicesdk

import (
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
)

type deprecated struct {
	container     container.Container
	keeperFactory KeeperFactory
}

func NewDeprecated(c container.Container, keeperFactory KeeperFactory) Service {
	return &deprecated{
		container:     c,
		keeperFactory: keeperFactory,
	}
}

// Create creates a new service from definition.
func (s *deprecated) Create(srv *service.Service) (*service.Service, error) {
	return create(nil, s.container, s.keeperFactory, srv)
}

// Delete deletes the service by hash.
func (s *deprecated) Delete(hash hash.Hash) error {
	keeper, err := s.keeperFactory(nil)
	if err != nil {
		return err
	}
	return keeper.Delete(hash)
}

// Get returns the service that matches given hash.
func (s *deprecated) Get(hash hash.Hash) (*service.Service, error) {
	keeper, err := s.keeperFactory(nil)
	if err != nil {
		return nil, err
	}
	return keeper.Get(hash)
}

// List returns all services.
func (s *deprecated) List() ([]*service.Service, error) {
	keeper, err := s.keeperFactory(nil)
	if err != nil {
		return nil, err
	}
	return keeper.All()
}
