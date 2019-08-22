package servicesdk

import (
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
)

type deprecated struct {
	container container.Container
	db        *database.ServiceDB
}

func NewDeprecated(c container.Container, db *database.ServiceDB) Service {
	return &deprecated{
		container: c,
		db:        db,
	}
}

// Create creates a new service from definition.
func (s *deprecated) Create(srv *service.Service) (*service.Service, error) {
	return create(s.container, s.db, srv)
}

// Delete deletes the service by hash.
func (s *deprecated) Delete(hash hash.Hash) error {
	return s.db.Delete(hash)
}

// Get returns the service that matches given hash.
func (s *deprecated) Get(hash hash.Hash) (*service.Service, error) {
	return s.db.Get(hash)
}

// List returns all services.
func (s *deprecated) List() ([]*service.Service, error) {
	return s.db.All()
}
