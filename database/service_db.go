package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/mesg-foundation/core/service"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const (
	aliasKeyPrefix = "alias_"
	idKeyPrefix    = "id_"
)

var (
	errCannotSaveWithoutID    = errors.New("database: can't save service without id")
	errCannotSaveWithoutAlias = errors.New("database: can't save service without alias")
	errAliasSameLen           = errors.New("database: service alias can't have same length as id")
)

// ServiceDB describes the API of database package.
type ServiceDB interface {
	// Save saves a service to database.
	Save(s *service.Service) error

	// Get gets a service from database by its unique id
	// or unique alias.
	Get(idOrAlias string) (*service.Service, error)

	// Delete deletes a service from database by its unique id
	// or unique alias.
	Delete(idOrAlias string) error

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
func (d *LevelDBServiceDB) unmarshal(id string, value []byte) (*service.Service, error) {
	var s service.Service
	if err := json.Unmarshal(value, &s); err != nil {
		return nil, &DecodeError{ID: id}
	}
	return &s, nil
}

// All returns every service in database.
func (d *LevelDBServiceDB) All() ([]*service.Service, error) {
	var (
		services []*service.Service
		iter     = d.db.NewIterator(util.BytesPrefix([]byte(idKeyPrefix)), nil)
	)
	for iter.Next() {
		id := strings.TrimPrefix(string(iter.Key()), idKeyPrefix)
		s, err := d.unmarshal(id, iter.Value())
		if err != nil {
			// NOTE: Ignore all decode errors (possibly due to a service
			// structure change or database corruption)
			if decodeErr, ok := err.(*DecodeError); ok {
				logrus.WithField("service", decodeErr.ID).Warning(decodeErr.Error())
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
func (d *LevelDBServiceDB) Delete(idOrAlias string) error {
	tx, err := d.db.OpenTransaction()
	if err != nil {
		return err
	}

	idBytes, err := d.db.Get([]byte(aliasKeyPrefix+idOrAlias), nil)
	if err != nil && err != leveldb.ErrNotFound {
		tx.Discard()
		return err
	}

	id := string(idBytes)
	alias := ""
	if id == "" {
		id = idOrAlias
	} else {
		alias = idOrAlias
	}

	batch := &leveldb.Batch{}
	batch.Delete([]byte(idKeyPrefix + id))
	if alias != "" {
		batch.Delete([]byte(aliasKeyPrefix + alias))
	}
	if err := tx.Write(batch, nil); err != nil {
		tx.Discard()
		return err
	}
	return tx.Commit()
}

// Get retrives service from database.
func (d *LevelDBServiceDB) Get(idOrAlias string) (*service.Service, error) {
	// check if key is an alias, if yes then save id.
	id, err := d.db.Get([]byte(aliasKeyPrefix+idOrAlias), nil)
	if err != nil && err != leveldb.ErrNotFound {
		return nil, err
	} else if err == nil {
		idOrAlias = string(id)
	}

	b, err := d.db.Get([]byte(idKeyPrefix+idOrAlias), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, &ErrNotFound{ID: idOrAlias}
		}
		return nil, err
	}
	return d.unmarshal(idOrAlias, b)
}

// Save stores service in database.
func (d *LevelDBServiceDB) Save(s *service.Service) error {
	// check service
	if s.ID == "" {
		return errCannotSaveWithoutID
	}
	if s.Alias == "" {
		return errCannotSaveWithoutAlias
	}
	if len(s.ID) == len(s.Alias) {
		return errAliasSameLen
	}

	// encode service
	b, err := d.marshal(s)
	if err != nil {
		return err
	}

	batch := &leveldb.Batch{}
	batch.Put([]byte(idKeyPrefix+s.ID), b)
	batch.Put([]byte(aliasKeyPrefix+s.Alias), []byte(s.ID))
	return d.db.Write(batch, nil)
}

// Close closes database.
func (d *LevelDBServiceDB) Close() error {
	return d.db.Close()
}

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
	return fmt.Sprintf("database: could not decode service %q", e.ID)
}

// IsErrNotFound returns true if err is type of ErrNotFound, false otherwise.
func IsErrNotFound(err error) bool {
	_, ok := err.(*ErrNotFound)
	return ok
}
