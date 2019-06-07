package database

 import (
	"encoding/json"

 	"github.com/mesg-foundation/core/instance"
	"github.com/syndtr/goleveldb/leveldb"
)

 // InstanceDB describes the API of Instance database.
type InstanceDB interface {
	// Get retrives instance by instance hash.
	Get(hash string) (*instance.Instance, error)

 	// Save saves instance to database.
	Save(i *instance.Instance) error

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
func (d *LevelDBInstanceDB) unmarshal(id string, value []byte) (*instance.Instance, error) {
	var s instance.Instance
	if err := json.Unmarshal(value, &s); err != nil {
		return nil, &DecodeError{ID: id}
	}
	return &s, nil
}

 // Get retrives instance by instance hash.
func (d *LevelDBInstanceDB) Get(hash string) (*instance.Instance, error) {
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

 // Save saves instance to database.
func (d *LevelDBInstanceDB) Save(i *instance.Instance) error {
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

 // Close closes database.
func (d *LevelDBInstanceDB) Close() error {
	return d.db.Close()
}
