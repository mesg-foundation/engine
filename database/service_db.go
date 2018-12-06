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
	sidKeyPrefix = "sid_"
)

var (
	errCannotSaveWithoutSID = errors.New("database: can't save service without sid")
)

// ServiceDB describes the API of database package.
type ServiceDB interface {
	// Save saves a service to database.
	Save(s *service.Service) error

	// Get gets a service from database by its unique id
	// or unique sid.
	Get(SID string) (*service.Service, error)

	// Delete deletes a service from database by its unique id
	// or unique sid.
	Delete(SID string) error

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
func (d *LevelDBServiceDB) unmarshal(sid string, value []byte) (*service.Service, error) {
	var s service.Service
	if err := json.Unmarshal(value, &s); err != nil {
		return nil, &DecodeError{SID: sid}
	}
	return &s, nil
}

// All returns every service in database.
func (d *LevelDBServiceDB) All() ([]*service.Service, error) {
	var (
		services []*service.Service
		iter     = d.db.NewIterator(util.BytesPrefix([]byte(sidKeyPrefix)), nil)
	)
	for iter.Next() {
		sid := strings.TrimPrefix(string(iter.Key()), sidKeyPrefix)
		s, err := d.unmarshal(sid, iter.Value())
		if err != nil {
			// NOTE: Ignore all decode errors (possibly due to a service
			// structure change or database corruption)
			if decodeErr, ok := err.(*DecodeError); ok {
				logrus.WithField("service", decodeErr.SID).Warning(decodeErr.Error())
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
func (d *LevelDBServiceDB) Delete(SID string) error {
	tx, err := d.db.OpenTransaction()
	if err != nil {
		return err
	}
	if _, err := d.get(tx, SID); err != nil {
		tx.Discard()
		return err
	}
	if err := tx.Delete([]byte(sidKeyPrefix+SID), nil); err != nil {
		tx.Discard()
		return err
	}
	return tx.Commit()
}

// Get retrives service from database.
func (d *LevelDBServiceDB) Get(SID string) (*service.Service, error) {
	return d.get(d.db, SID)
}

// get retrives service from database by using r reader.
func (d *LevelDBServiceDB) get(r leveldb.Reader, SID string) (*service.Service, error) {
	// get the service
	b, err := r.Get([]byte(sidKeyPrefix+SID), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, &ErrNotFound{SID: SID}
		}
		return nil, err
	}
	return d.unmarshal(SID, b)
}

// Save stores service in database.
// If there is an another service that uses the same sid, it'll be deleted.
func (d *LevelDBServiceDB) Save(s *service.Service) error {
	// check service
	if s.SID == "" {
		return errCannotSaveWithoutSID
	}

	// encode service
	b, err := d.marshal(s)
	if err != nil {
		return err
	}

	// save service.
	return d.db.Put([]byte(sidKeyPrefix+s.SID), b, nil)
}

// Close closes database.
func (d *LevelDBServiceDB) Close() error {
	return d.db.Close()
}

// ErrNotFound is an not found error.
type ErrNotFound struct {
	SID string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("database: service %s not found", e.SID)
}

// DecodeError represents a service impossible to decode.
type DecodeError struct {
	SID string
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("database: could not decode service %q", e.SID)
}

// IsErrNotFound returns true if err is type of ErrNotFound, false otherwise.
func IsErrNotFound(err error) bool {
	_, ok := err.(*ErrNotFound)
	return ok
}
