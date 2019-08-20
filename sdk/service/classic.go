package servicesdk

import (
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
)

type Classic struct {
	logic *logic
}

func NewClassic(c container.Container, keeperFactory KeeperFactory) *Classic {
	return &Classic{
		logic: newLogic(c, keeperFactory),
	}
}

// Create creates a new service from definition.
func (s *Classic) Create(srv *service.Service) (*service.Service, error) {
	return s.logic.create(nil, srv)
}

// Delete deletes the service by hash.
func (s *Classic) Delete(hash hash.Hash) error {
	return s.logic.delete(nil, hash)
}

// Get returns the service that matches given hash.
func (s *Classic) Get(hash hash.Hash) (*service.Service, error) {
	return s.logic.get(nil, hash)
}

// List returns all services.
func (s *Classic) List() ([]*service.Service, error) {
	return s.logic.list(nil)
}
