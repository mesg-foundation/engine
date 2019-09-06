package database

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/mesg-foundation/engine/database/store"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
	"github.com/sirupsen/logrus"
)

var (
	errCannotSaveWithoutHash = errors.New("database: can't save service without hash")
)

// ServiceDB is a database for storing service definition.
type ServiceDB struct {
	s   store.Store
	cdc *codec.Codec
}

// NewServiceDB returns the database which is located under given path.
func NewServiceDB(s store.Store, cdc *codec.Codec) *ServiceDB {
	return &ServiceDB{
		s:   s,
		cdc: cdc,
	}
}

// unmarshal returns the service from byte slice.
func (d *ServiceDB) unmarshalService(hash hash.Hash, value []byte) (*service.Service, error) {
	var s service.Service
	if err := d.cdc.UnmarshalBinaryBare(value, &s); err != nil {
		return nil, fmt.Errorf("database: could not decode service %q: %s", hash, err)
	}
	return &s, nil
}

// Exist check if service with given hash exist.
func (d *ServiceDB) Exist(hash hash.Hash) (bool, error) {
	return d.s.Has(hash)
}

// All returns every service in database.
func (d *ServiceDB) All() ([]*service.Service, error) {
	var (
		services []*service.Service
		iter     = d.s.NewIterator()
	)
	for iter.Next() {
		hash := hash.Hash(iter.Key())
		s, err := d.unmarshalService(hash, iter.Value())
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
func (d *ServiceDB) Delete(hash hash.Hash) error {
	return d.s.Delete(hash)
}

// Get retrives service from database.
func (d *ServiceDB) Get(hash hash.Hash) (*service.Service, error) {
	b, err := d.s.Get(hash)
	if err != nil {
		return nil, err
	}
	return d.unmarshalService(hash, b)
}

// Save stores service in database.
// If there is an another service that uses the same sid, it'll be deleted.
func (d *ServiceDB) Save(s *service.Service) error {
	if len(s.Hash) == 0 {
		return errCannotSaveWithoutHash
	}
	b, err := d.cdc.MarshalBinaryBare(s)
	if err != nil {
		return err
	}
	return d.s.Put(s.Hash, b)
}

// Close closes database.
func (d *ServiceDB) Close() error {
	return d.s.Close()
}
