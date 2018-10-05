package database

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mesg-foundation/core/service"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

// ErrNotFound is an not found error.
type ErrNotFound struct {
	ID string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("database: service %s not found", e.ID)
}

// DecodeError represents a service impossible to decode.
type DecodeError struct {
	ID string
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("Database services: Could not decode service %q", e.ID)
}

// IsErrNotFound returs true if err is type of ErrNotFound, false otherwise.
func IsErrNotFound(err error) bool {
	_, ok := err.(*ErrNotFound)
	return ok
}

// ServiceDB is a database for storing service definition.
type ServiceDB struct {
	db *leveldb.DB
}

// NewServiceDB returns database which is located under given path.
func NewServiceDB(path string) (*ServiceDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &ServiceDB{db: db}, nil
}

// marshal returns the byte slice from service.
func (db *ServiceDB) marshal(s *service.Service) ([]byte, error) {
	return json.Marshal(s)
}

// unmarshal returns the service from byte slice.
func (db *ServiceDB) unmarshal(id string, value []byte) (*service.Service, error) {
	var s service.Service
	if err := json.Unmarshal(value, &s); err != nil {
		return nil, &DecodeError{ID: id}
	}
	return &s, nil
}

// All returns every service in database.
func (db *ServiceDB) All() ([]*service.Service, error) {
	var (
		services []*service.Service
		iter     = db.db.NewIterator(nil, nil)
	)

	for iter.Next() {
		s, err := db.unmarshal(string(iter.Key()), iter.Value())
		if err != nil {
			// NOTE: Ignore all decode errors (possibly due to a service
			// structure change or database corruption)
			if decodeErr, ok := err.(*DecodeError); ok {
				logrus.WithField("service", decodeErr.ID).Warning(decodeErr.Error())
				continue
			}
			return nil, err
		}
		services = append(services, s)
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return nil, err
	}

	return services, nil
}

// Delete deletes service from database.
func (db *ServiceDB) Delete(id string) error {
	return db.db.Delete([]byte(id), nil)
}

// Get retrives service from database.
func (db *ServiceDB) Get(id string) (*service.Service, error) {
	b, err := db.db.Get([]byte(id), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, &ErrNotFound{ID: id}
		}
		return nil, err
	}

	return db.unmarshal(id, b)
}

// Save stores service in database.
func (db *ServiceDB) Save(s *service.Service) error {
	if s.ID == "" {
		return errors.New("database: can't save service without id")
	}
	b, err := db.marshal(s)
	if err != nil {
		return err
	}

	return db.db.Put([]byte(s.ID), b, nil)
}

// Close closes database.
func (db *ServiceDB) Close() error {
	return db.db.Close()
}
