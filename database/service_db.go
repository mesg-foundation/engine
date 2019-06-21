package database

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mesg-foundation/core/hash"
	"github.com/mesg-foundation/core/service"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	errSaveInstanceWithoutHash = errors.New("database: can't save instance without hash")
	errCannotSaveWithoutHash   = errors.New("database: can't save service without hash")
)

// ServiceDB describes the API of database package.
type ServiceDB interface {
	// Save saves a service to database.
	Save(s *service.Service) error

	// Get gets a service from database by its unique hash.
	Get(hash hash.Hash) (*service.Service, error)

	// Delete deletes a service from database by its unique hash.
	Delete(hash hash.Hash) error

	// All returns all services from database.
	All() ([]*service.Service, error)

	// Close closes underlying database connection.
	Close() error
}

// LevelDBServiceDB is a database for storing service definition.
type LevelDBServiceDB struct {
	db *leveldb.DB
}

// NewServiceDB returns the database which is located under given path.
func NewServiceDB(path string) (*LevelDBServiceDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDBServiceDB{db: db}, nil
}

// marshal returns the byte slice from service.
func (d *LevelDBServiceDB) marshal(s *service.Service) ([]byte, error) {
	return json.Marshal(s)
}

// unmarshal returns the service from byte slice.
func (d *LevelDBServiceDB) unmarshal(hash, value []byte) (*service.Service, error) {
	var s service.Service
	if err := json.Unmarshal(value, &s); err != nil {
		return nil, &DecodeError{hash: hash}
	}
	return &s, nil
}

// All returns every service in database.
func (d *LevelDBServiceDB) All() ([]*service.Service, error) {
	var (
		services []*service.Service
		iter     = d.db.NewIterator(nil, nil)
	)
	for iter.Next() {
		s, err := d.unmarshal(iter.Key(), iter.Value())
		if err != nil {
			// NOTE: Ignore all decode errors (possibly due to a service
			// structure change or database corruption)
			if decodeErr, ok := err.(*DecodeError); ok {
				logrus.WithField("service", decodeErr.hash.String()).Warning(decodeErr.Error())
				continue
			}
			iter.Release()
			return nil, err
		}
		services = append(services, s)
	}
	iter.Release()
	return services, iter.Error()
}

// Delete deletes service from database.
func (d *LevelDBServiceDB) Delete(hash hash.Hash) error {
	tx, err := d.db.OpenTransaction()
	if err != nil {
		return err
	}
	if _, err := tx.Get(hash, nil); err != nil {
		tx.Discard()
		if err == leveldb.ErrNotFound {
			return &ErrNotFound{hash: hash}
		}
		return err
	}
	if err := tx.Delete(hash, nil); err != nil {
		tx.Discard()
		return err
	}
	return tx.Commit()

}

// Get retrives service from database.
func (d *LevelDBServiceDB) Get(hash hash.Hash) (*service.Service, error) {
	b, err := d.db.Get(hash, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, &ErrNotFound{hash: hash}
		}
		return nil, err
	}
	return d.unmarshal(hash, b)
}

// Save stores service in database.
// If there is an another service that uses the same sid, it'll be deleted.
func (d *LevelDBServiceDB) Save(s *service.Service) error {
	if s.Hash.IsZero() {
		return errCannotSaveWithoutHash
	}

	b, err := d.marshal(s)
	if err != nil {
		return err
	}
	return d.db.Put(s.Hash, b, nil)
}

// Close closes database.
func (d *LevelDBServiceDB) Close() error {
	return d.db.Close()
}

// ErrNotFound is an not found error.
type ErrNotFound struct {
	hash hash.Hash
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("database: service %q not found", e.hash)
}

// DecodeError represents a service impossible to decode.
type DecodeError struct {
	hash hash.Hash
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("database: could not decode service %q", e.hash)
}

// IsErrNotFound returns true if err is type of ErrNotFound, false otherwise.
func IsErrNotFound(err error) bool {
	_, ok := err.(*ErrNotFound)
	return ok
}
