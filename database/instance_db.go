package database

import (
	"encoding/json"

	"github.com/mesg-foundation/core/hash"
	"github.com/mesg-foundation/core/instance"
	"github.com/syndtr/goleveldb/leveldb"
)

// InstanceDB describes the API of Instance database.
type InstanceDB interface {
	// Get retrives instance by instance hash.
	Get(hash hash.Hash) (*instance.Instance, error)

	// GetAll retrieves all instances.
	GetAll() ([]*instance.Instance, error)

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
		return nil, &DecodeError{hash: hash}
	}
	return &s, nil
}

// Get retrives instance by instance hash.
func (d *LevelDBInstanceDB) Get(hash hash.Hash) (*instance.Instance, error) {
	b, err := d.db.Get(hash, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, &ErrNotFound{hash: hash}
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

// Delete deletes service from database.
func (d *LevelDBInstanceDB) Delete(hash hash.Hash) error {
	return d.db.Delete(hash, nil)
}
