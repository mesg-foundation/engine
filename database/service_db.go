package database

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/keeper"
	"github.com/mesg-foundation/engine/service"
	"github.com/sirupsen/logrus"
)

var (
	errCannotSaveWithoutHash = errors.New("database: can't save service without hash")
)

// ServiceKeeper is a database for storing service definition.
type ServiceKeeper struct {
}

// NewServiceKeeper returns the database which is located under given path.
func NewServiceKeeper() (*ServiceKeeper, error) {
	return &ServiceKeeper{}, nil
}

// marshal returns the byte slice from service.
func (k *ServiceKeeper) marshal(s *service.Service) ([]byte, error) {
	return json.Marshal(s)
}

// unmarshal returns the service from byte slice.
func (k *ServiceKeeper) unmarshal(hash hash.Hash, value []byte) (*service.Service, error) {
	var s service.Service
	if err := json.Unmarshal(value, &s); err != nil {
		return nil, fmt.Errorf("database: could not decode service %q: %s", hash, err)
	}
	return &s, nil
}

// All returns every service in database.
func (k *ServiceKeeper) All(s keeper.Store) ([]*service.Service, error) {
	var (
		services []*service.Service
		iter     = s.NewIterator()
	)
	for iter.Next() {
		hash := hash.Hash(iter.Key())
		s, err := k.unmarshal(hash, iter.Value())
		if err != nil {
			// NOTE: Ignore all decode errors (possibly due to a service
			// structure change or database corruption)
			logrus.WithField("service", hash.String()).Warning(err.Error())
			continue
		}
		services = append(services, s)
	}
	iter.Release()
	return services, iter.Error()
}

// Delete deletes service from database.
func (k *ServiceKeeper) Delete(s keeper.Store, hash hash.Hash) error {
	has, err := s.Has(hash)
	if err != nil {
		return err
	}
	if !has {
		return &ErrNotFound{resource: "service", hash: hash}
	}
	return s.Delete(hash)
}

// Get retrives service from database.
func (k *ServiceKeeper) Get(s keeper.Store, hash hash.Hash) (*service.Service, error) {
	has, err := s.Has(hash)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, &ErrNotFound{resource: "service", hash: hash}
	}
	b, err := s.Get(hash)
	if err != nil {
		return nil, err
	}
	return k.unmarshal(hash, b)
}

// Save stores service in database.
// If there is an another service that uses the same sid, it'll be deleted.
func (k *ServiceKeeper) Save(s keeper.Store, srv *service.Service) error {
	if srv.Hash.IsZero() {
		return errCannotSaveWithoutHash
	}
	b, err := k.marshal(srv)
	if err != nil {
		return err
	}
	return s.Put(srv.Hash, b)
}

// Close closes database.
func (k *ServiceKeeper) Close(s keeper.Store) error {
	return s.Close()
}

// ErrNotFound is an not found error.
type ErrNotFound struct {
	hash     hash.Hash
	resource string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("database: %s %q not found", e.resource, e.hash)
}

// IsErrNotFound returns true if err is type of ErrNotFound, false otherwise.
func IsErrNotFound(err error) bool {
	_, ok := err.(*ErrNotFound)
	return ok
}
