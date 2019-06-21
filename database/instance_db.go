package database

import (
	"bytes"
	"encoding/json"

	"github.com/mesg-foundation/core/instance"
	"github.com/mr-tron/base58"
	"github.com/syndtr/goleveldb/leveldb"
)

// InstanceDB describes the API of Instance database.
type InstanceDB interface {
	// Get retrives instance by instance hash.
	Get(hash []byte) (*instance.Instance, error)

	// GetAll retrieves all instances.
	GetAll() ([]*instance.Instance, error)

	// GetAllByService retrieves all instances of service by service's hash.
	GetAllByService(serviceHash []byte) ([]*instance.Instance, error)

	// Save saves instance to database.
	Save(i *instance.Instance) error

	// Delete an instance by instance hash.
	Delete(hash []byte) error

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
		return nil, &DecodeError{ID: base58.Encode(hash)}
	}
	return &s, nil
}

// Get retrives instance by instance hash.
func (d *LevelDBInstanceDB) Get(hash []byte) (*instance.Instance, error) {
	b, err := d.db.Get(hash, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, &ErrNotFound{ID: base58.Encode(hash)}
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
func (d *LevelDBInstanceDB) GetAllByService(serviceHash []byte) ([]*instance.Instance, error) {
	instances, err := d.GetAll()
	if err != nil {
		return nil, err
	}
	someInstances := []*instance.Instance{}
	for _, instance := range instances {
		if bytes.Equal(instance.ServiceHash, serviceHash) {
			someInstances = append(someInstances, instance)
		}
	}
	return someInstances, nil
}

// Save saves instance to database.
func (d *LevelDBInstanceDB) Save(i *instance.Instance) error {
	if len(i.Hash) == 0 {
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

// Delete deletes service from database.
func (d *LevelDBInstanceDB) Delete(hash []byte) error {
	return d.db.Delete(hash, nil)
}
