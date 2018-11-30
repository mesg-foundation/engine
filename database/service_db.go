package database

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mesg-foundation/core/service"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

const aliasesPathSuffix = "_aliases"

var (
	errCannotSaveWithoutID    = errors.New("database: can't save service without id")
	errCannotSaveWithoutAlias = errors.New("database: can't save service without alias")
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
	db      *leveldb.DB
	aliases *leveldb.DB
}

// NewServiceDB returns the database which is located under given path.
func NewServiceDB(path string) (*LevelDBServiceDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	aliases, err := leveldb.OpenFile(path+aliasesPathSuffix, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDBServiceDB{
		db:      db,
		aliases: aliases,
	}, nil
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
		iter     = d.db.NewIterator(nil, nil)
	)
	for iter.Next() {
		s, err := d.unmarshal(string(iter.Key()), iter.Value())
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
	id, alias, err := d.getIDAndAlias(idOrAlias)
	if err != nil {
		return err
	}
	// delete alias
	if alias != "" {
		if err := d.aliases.Delete([]byte(alias), nil); err != nil {
			return err
		}
	}
	// delete service
	return d.db.Delete([]byte(id), nil)
}

// Get retrives service from database.
func (d *LevelDBServiceDB) Get(idOrAlias string) (*service.Service, error) {
	id, _, err := d.getIDAndAlias(idOrAlias)
	if err != nil {
		return nil, err
	}
	b, err := d.db.Get([]byte(id), nil)
	if err != nil {
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
	// encode service
	b, err := d.marshal(s)
	if err != nil {
		return err
	}
	// save service
	if err := d.db.Put([]byte(s.ID), b, nil); err != nil {
		return err
	}
	// save alias if exist
	if s.Alias == "" {
		return errCannotSaveWithoutAlias
	}
	if err := d.aliases.Put([]byte(s.Alias), []byte(s.ID), nil); err != nil {
		return err
	}
	return nil
}

// Close closes database.
func (d *LevelDBServiceDB) Close() error {
	if err := d.db.Close(); err != nil {
		return err
	}
	return d.aliases.Close()
}

// getIDAndAlias returns and separates id from alias.
// id will be filled if there are no errors.
// returned alias can be empty.
func (d *LevelDBServiceDB) getIDAndAlias(idOrAlias string) (id string, alias string, err error) {
	// check if idOrAlias is an id
	// if it is already an id, return it with correspond alias if exist
	isID, err := d.db.Has([]byte(idOrAlias), nil)
	if err != nil {
		return "", "", err
	}
	if isID {
		id = idOrAlias
		alias, err = d.getAliasFromID(id)
		if err != nil {
			return "", "", err
		}
		return id, alias, nil
	}
	// check if idOrAlias is an alias
	// if it is an alias, return it with the corresponding id
	alias = idOrAlias
	idB, err := d.aliases.Get([]byte(alias), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return "", "", &ErrNotFound{ID: alias}
		}
		return "", "", err
	}
	id = string(idB)
	return id, alias, nil
}

// getAliasFromID returns the alias pointing to the id if exist.
func (d *LevelDBServiceDB) getAliasFromID(id string) (alias string, err error) {
	iter := d.aliases.NewIterator(nil, nil)
	for iter.Next() {
		if string(iter.Value()) == id {
			alias = string(iter.Key())
			break
		}
	}
	iter.Release()
	return alias, iter.Error()
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
	return fmt.Sprintf("database: Could not decode service %q", e.ID)
}

// IsErrNotFound returns true if err is type of ErrNotFound, false otherwise.
func IsErrNotFound(err error) bool {
	_, ok := err.(*ErrNotFound)
	return ok
}
