package database

import (
	"errors"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	errCannotSaveProcessWithoutHash = errors.New("database: can't save process without hash")
)

// ProcessDB describes the API of database package.
type ProcessDB interface {
	// Save saves a process to database.
	Save(s *process.Process) error

	// Get gets a process from database by its unique hash.
	Get(hash hash.Hash) (*process.Process, error)

	// Delete deletes a process from database by its unique hash.
	Delete(hash hash.Hash) error

	// All returns all processes from database.
	All() ([]*process.Process, error)

	// Close closes underlying database connection.
	Close() error
}

// LevelDBProcessDB is a database for storing processes definition.
type LevelDBProcessDB struct {
	db *leveldb.DB
}

// NewProcessDB returns the database which is located under given path.
func NewProcessDB(path string) (*LevelDBProcessDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDBProcessDB{db: db}, nil
}

// marshal returns the byte slice from process.
func (d *LevelDBProcessDB) marshal(s proto.Message) ([]byte, error) {
	return proto.Marshal(s)
}

// unmarshal returns the process from byte slice.
func (d *LevelDBProcessDB) unmarshal(hash hash.Hash, value []byte) (*process.Process, error) {
	var s process.Process
	if err := proto.Unmarshal(value, &s); err != nil {
		return nil, fmt.Errorf("database: could not decode process %q: %s", hash, err)
	}
	return &s, nil
}

// All returns every process in database.
func (d *LevelDBProcessDB) All() ([]*process.Process, error) {
	var (
		processes []*process.Process
		iter      = d.db.NewIterator(nil, nil)
	)
	for iter.Next() {
		hash := hash.Hash(iter.Key())
		s, err := d.unmarshal(hash, iter.Value())
		if err != nil {
			// NOTE: Ignore all decode errors (possibly due to a process
			// structure change or database corruption)
			logrus.WithField("process", hash.String()).Warning(err.Error())
			continue
		}
		processes = append(processes, s)
	}
	iter.Release()
	return processes, iter.Error()
}

// Delete deletes process from database.
func (d *LevelDBProcessDB) Delete(hash hash.Hash) error {
	tx, err := d.db.OpenTransaction()
	if err != nil {
		return err
	}
	if _, err := tx.Get(hash, nil); err != nil {
		tx.Discard()
		return err
	}
	if err := tx.Delete(hash, nil); err != nil {
		tx.Discard()
		return err
	}
	return tx.Commit()
}

// Get retrives process from database.
func (d *LevelDBProcessDB) Get(hash hash.Hash) (*process.Process, error) {
	b, err := d.db.Get(hash, nil)
	if err != nil {
		return nil, err
	}
	return d.unmarshal(hash, b)
}

// Save stores process in database.
// If there is an another process that uses the same sid, it'll be deleted.
func (d *LevelDBProcessDB) Save(s *process.Process) error {
	if s.Hash.IsZero() {
		return errCannotSaveProcessWithoutHash
	}

	b, err := d.marshal(s)
	if err != nil {
		return err
	}
	return d.db.Put(s.Hash, b, nil)
}

// Close closes database.
func (d *LevelDBProcessDB) Close() error {
	return d.db.Close()
}
