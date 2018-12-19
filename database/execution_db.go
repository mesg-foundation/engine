package database

import (
	"encoding/json"
	"errors"

	"github.com/mesg-foundation/core/execution"
	"github.com/syndtr/goleveldb/leveldb"
)

// ExecutionDB exposes all the functionalities
type ExecutionDB interface {
	Find(executionID string) (*execution.Execution, error)
	Save(execution *execution.Execution) error
	Close() error
	OpenTransaction() (*leveldb.Transaction, error)
	FindWithTx(tx leveldb.Reader, executionID string) (*execution.Execution, error)
	SaveWithTx(tx *leveldb.Transaction, execution *execution.Execution) error
}

// LevelDBExecutionDB is a concrete implementation of the DB interface
type LevelDBExecutionDB struct {
	db *leveldb.DB
}

// NewExecutionDB creates a new DB instance
func NewExecutionDB(path string) (*LevelDBExecutionDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	return &LevelDBExecutionDB{db: db}, nil
}

// Find the execution based on an executionID, returns an error if not found
func (db *LevelDBExecutionDB) Find(executionID string) (*execution.Execution, error) {
	return db.FindWithTx(db.db, executionID)
}

// FindWithTx the execution based on an executionID, returns an error if not found
func (db *LevelDBExecutionDB) FindWithTx(tx leveldb.Reader, executionID string) (*execution.Execution, error) {
	data, err := tx.Get([]byte(executionID), nil)
	if err != nil {
		return nil, err
	}
	var execution execution.Execution
	if err := json.Unmarshal(data, &execution); err != nil {
		return nil, err
	}
	return &execution, nil
}

// Save an instance of executable in the database
// Returns an error if anything from marshaling to database saving goes wrong
func (db *LevelDBExecutionDB) Save(execution *execution.Execution) error {
	tx, err := db.db.OpenTransaction()
	if err != nil {
		return err
	}
	if err := db.SaveWithTx(tx, execution); err != nil {
		tx.Discard()
		return err
	}
	return tx.Commit()
}

// SaveWithTx an instance of executable in the database
// Returns an error if anything from marshaling to database saving goes wrong
func (db *LevelDBExecutionDB) SaveWithTx(tx *leveldb.Transaction, execution *execution.Execution) error {
	if execution.ID == "" {
		return errors.New("database: can't save service without id")
	}
	data, err := json.Marshal(execution)
	if err != nil {
		return err
	}
	return tx.Put([]byte(execution.ID), data, nil)
}

// Close closes database.
func (db *LevelDBExecutionDB) Close() error {
	return db.db.Close()
}

// OpenTransaction creates an new transaction on this db.
func (db *LevelDBExecutionDB) OpenTransaction() (*leveldb.Transaction, error) {
	return db.db.OpenTransaction()
}
