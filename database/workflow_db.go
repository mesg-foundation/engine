package database

import (
	"encoding/json"
	"fmt"

	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/workflow"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
)

// WorkflowDB describes the API of database package.
type WorkflowDB interface {
	// Save saves a workflow to database.
	Save(s *workflow.Workflow) error

	// Get gets a workflow from database by its unique hash.
	Get(hash hash.Hash) (*workflow.Workflow, error)

	// Delete deletes a workflow from database by its unique hash.
	Delete(hash hash.Hash) error

	// All returns all workflows from database.
	All() ([]*workflow.Workflow, error)

	// Close closes underlying database connection.
	Close() error
}

// LevelDBWorkflowDB is a database for storing workflows definition.
type LevelDBWorkflowDB struct {
	db *leveldb.DB
}

// NewWorkflowDB returns the database which is located under given path.
func NewWorkflowDB(path string) (*LevelDBWorkflowDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDBWorkflowDB{db: db}, nil
}

// marshal returns the byte slice from workflow.
func (d *LevelDBWorkflowDB) marshal(s *workflow.Workflow) ([]byte, error) {
	return json.Marshal(s)
}

// unmarshal returns the workflow from byte slice.
func (d *LevelDBWorkflowDB) unmarshal(hash hash.Hash, value []byte) (*workflow.Workflow, error) {
	var s workflow.Workflow
	if err := json.Unmarshal(value, &s); err != nil {
		return nil, fmt.Errorf("database: could not decode workflow %q: %s", hash, err)
	}
	return &s, nil
}

// All returns every workflow in database.
func (d *LevelDBWorkflowDB) All() ([]*workflow.Workflow, error) {
	var (
		workflows []*workflow.Workflow
		iter      = d.db.NewIterator(nil, nil)
	)
	for iter.Next() {
		hash := hash.Hash(iter.Key())
		s, err := d.unmarshal(hash, iter.Value())
		if err != nil {
			// NOTE: Ignore all decode errors (possibly due to a workflow
			// structure change or database corruption)
			logrus.WithField("workflow", hash.String()).Warning(err.Error())
			continue
		}
		workflows = append(workflows, s)
	}
	iter.Release()
	return workflows, iter.Error()
}

// Delete deletes workflow from database.
func (d *LevelDBWorkflowDB) Delete(hash hash.Hash) error {
	tx, err := d.db.OpenTransaction()
	if err != nil {
		return err
	}
	if _, err := tx.Get(hash, nil); err != nil {
		tx.Discard()
		if err == leveldb.ErrNotFound {
			return &ErrNotFound{resource: "workflow", hash: hash}
		}
		return err
	}
	if err := tx.Delete(hash, nil); err != nil {
		tx.Discard()
		return err
	}
	return tx.Commit()

}

// Get retrives workflow from database.
func (d *LevelDBWorkflowDB) Get(hash hash.Hash) (*workflow.Workflow, error) {
	b, err := d.db.Get(hash, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, &ErrNotFound{resource: "workflow", hash: hash}
		}
		return nil, err
	}
	return d.unmarshal(hash, b)
}

// Save stores workflow in database.
// If there is an another workflow that uses the same sid, it'll be deleted.
func (d *LevelDBWorkflowDB) Save(s *workflow.Workflow) error {
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
func (d *LevelDBWorkflowDB) Close() error {
	return d.db.Close()
}
