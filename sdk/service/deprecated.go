package servicesdk

import (
	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/database"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
)

// deprecated exposes service APIs of MESG.
type deprecated struct {
	container container.Container
	serviceDB *database.ServiceDB
}

// NewDeprecated creates a new Service SDK with given options.
func NewDeprecated(c container.Container, serviceDB *database.ServiceDB) Service {
	return &deprecated{
		container: c,
		serviceDB: serviceDB,
	}
}

// Create creates a new service from definition.
func (s *deprecated) Create(srv *service.Service) (*service.Service, error) {
	return create(s.container, s.serviceDB, srv)
}

// Delete deletes the service by hash.
func (s *deprecated) Delete(hash hash.Hash) error {
	return s.serviceDB.Delete(hash)
}

// Get returns the service that matches given hash.
func (s *deprecated) Get(hash hash.Hash) (*service.Service, error) {
	return s.serviceDB.Get(hash)
}

// List returns all services.
func (s *deprecated) List() ([]*service.Service, error) {
	return s.serviceDB.All()
}
