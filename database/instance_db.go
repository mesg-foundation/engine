package database

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mesg-foundation/core/hash"
	"github.com/mesg-foundation/core/instance"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	errSaveInstanceWithoutHash = errors.New("database: can't save instance without hash")
)

// InstanceDB describes the API of Instance database.
type InstanceDB interface {
	// Get retrives instance by instance hash.
	Get(hash hash.Hash) (*instance.Instance, error)

	// GetAll retrieves all instances.
	GetAll() ([]*instance.Instance, error)

	// GetAllByService retrieves all instances of service by service's hash.
	GetAllByService(serviceHash hash.Hash) ([]*instance.Instance, error)

	// Save saves instance to database.
	Save(i *instance.Instance) error

	// Delete an instance by instance hash.
	Delete(hash hash.Hash) error

	// Close closes underlying database connection.
	Close() error
}

// LevelDBInstanceDB is a database for storing services' instances.
type LevelDBInstanceDB struct {
	db *leveldb.DB
}

// NewInstanceDB returns the database which is located under given path.
func NewInstanceDB(path string) (*LevelDBInstanceDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDBInstanceDB{db: db}, nil
}

// marshal returns the byte slice from service.
func (d *LevelDBInstanceDB) marshal(i *instance.Instance) ([]byte, error) {
	return json.Marshal(i)
}

// unmarshal returns the service from byte slice.
func (d *LevelDBInstanceDB) unmarshal(hash, value []byte) (*instance.Instance, error) {
	var s instance.Instance
	if err := json.Unmarshal(value, &s); err != nil {
		return nil, &InstanceDecodeError{hash: hash}
	}
	return &s, nil
}

// Get retrives instance by instance hash.
func (d *LevelDBInstanceDB) Get(hash hash.Hash) (*instance.Instance, error) {
	b, err := d.db.Get(hash, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, &ErrInstanceNotFound{hash: hash}
		}
		return nil, err
	}
	return d.unmarshal(hash, b)
}

// GetAll retrieves all instances.
func (d *LevelDBInstanceDB) GetAll() ([]*instance.Instance, error) {
	instances := []*instance.Instance{}
	iter := d.db.NewIterator(nil, nil)
	for iter.Next() {
		i, err := d.unmarshal(iter.Key(), iter.Value())
		if err != nil {
			iter.Release()
			return nil, err
		}
		instances = append(instances, i)
	}
	iter.Release()
	return instances, iter.Error()
}

// GetAllByService retrieves all instances of service by service's hash.
func (d *LevelDBInstanceDB) GetAllByService(serviceHash hash.Hash) ([]*instance.Instance, error) {
	instances, err := d.GetAll()
	if err != nil {
		return nil, err
	}
	someInstances := []*instance.Instance{}
	for _, instance := range instances {
		if instance.ServiceHash.Equal(serviceHash) {
			someInstances = append(someInstances, instance)
		}
	}
	return someInstances, nil
}

// Save saves instance to database.
func (d *LevelDBInstanceDB) Save(i *instance.Instance) error {
	if i.Hash.IsZero() {
		return errSaveInstanceWithoutHash
	}

	// encode service
	b, err := d.marshal(i)
	if err != nil {
		return err
	}

	return d.db.Put(i.Hash, b, nil)
}

// Close closes database.
func (d *LevelDBInstanceDB) Close() error {
	return d.db.Close()
}

// Delete deletes an instance from database.
func (d *LevelDBInstanceDB) Delete(hash hash.Hash) error {
	return d.db.Delete(hash, nil)
}

// ErrInstanceNotFound is an not found error.
type ErrInstanceNotFound struct {
	hash hash.Hash
}

func (e *ErrInstanceNotFound) Error() string {
	return fmt.Sprintf("database: instance %q not found", e.hash)
}

// InstanceDecodeError represents a resource impossible to decode.
type InstanceDecodeError struct {
	hash hash.Hash
}

func (e *InstanceDecodeError) Error() string {
	return fmt.Sprintf("database: could not instance service %q", e.hash)
}

// IsErrInstanceNotFound returns true if err is type of ErrInstanceNotFound, false otherwise.
func IsErrInstanceNotFound(err error) bool {
	_, ok := err.(*ErrInstanceNotFound)
	return ok
}
