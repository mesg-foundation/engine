package database

import (
	"encoding/json"

	"github.com/mesg-foundation/core/service"
	"github.com/syndtr/goleveldb/leveldb"
)

// InstanceDB describes the API of Instance database.
type InstanceDB interface {
	// Get retrives instance by instance hash.
	Get(hash string) (*service.Instance, error)

	// GetByService retrives instance by service hash.
	GetByService(hash string) (*service.Instance, error)

	// Save saves instance to database.
	Save(i *service.Instance) error

	// Delete deletes instance from database by hash.
	Delete(hash string) error

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
func (d *LevelDBInstanceDB) marshal(i *service.Instance) ([]byte, error) {
	return json.Marshal(i)
}

// unmarshal returns the service from byte slice.
func (d *LevelDBInstanceDB) unmarshal(id string, value []byte) (*service.Instance, error) {
	var s service.Instance
	if err := json.Unmarshal(value, &s); err != nil {
		return nil, &DecodeError{ID: id}
	}
	return &s, nil
}

// Get retrives instance by instance hash.
func (d *LevelDBInstanceDB) Get(hash string) (*service.Instance, error) {
	tx, err := d.db.OpenTransaction()
	if err != nil {
		return nil, err
	}
	b, err := tx.Get([]byte(hash), nil)
	if err != nil {
		tx.Discard()
		if err == leveldb.ErrNotFound {
			return nil, &ErrNotFound{ID: hash}
		}
		return nil, err
	}
	i, err := d.unmarshal(hash, b)
	if err != nil {
		tx.Discard()
		return nil, err
	}
	return i, tx.Commit()
}

// GetByService retrives instance by its service's hash.
func (d *LevelDBInstanceDB) GetByService(hash string) (*service.Instance, error) {
	iter := d.db.NewIterator(nil, nil)
	for iter.Next() {
		instanceHash := string(iter.Key())
		i, err := d.unmarshal(instanceHash, iter.Value())
		if err != nil {
			iter.Release()
			return nil, err
		}
		if i.ServiceHash == hash {
			iter.Release()
			return i, nil
		}
	}
	iter.Release()
	if iter.Error() != nil {
		return nil, iter.Error()
	}
	return nil, &ErrNotFound{ID: hash}
}

// Save saves instance to database.
func (d *LevelDBInstanceDB) Save(i *service.Instance) error {
	// check service
	if i.Hash == "" {
		return errCannotSaveWithoutHash
	}

	// open database transaction
	tx, err := d.db.OpenTransaction()
	if err != nil {
		return err
	}

	// encode service
	b, err := d.marshal(i)
	if err != nil {
		tx.Discard()
		return err
	}

	// save instance with hash.
	if err := tx.Put([]byte(i.Hash), b, nil); err != nil {
		tx.Discard()
		return err
	}

	return tx.Commit()
}

// Delete deletes instance from database by hash.
func (d *LevelDBInstanceDB) Delete(hash string) error {
	tx, err := d.db.OpenTransaction()
	if err != nil {
		return err
	}
	if err := tx.Delete([]byte(hash), nil); err != nil {
		tx.Discard()
		return err
	}
	return tx.Commit()
}

// Close closes database.
func (d *LevelDBInstanceDB) Close() error {
	return d.db.Close()
}
