package database

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/mesg-foundation/core/execution"
	"github.com/syndtr/goleveldb/leveldb"
)

// ExecutionDB exposes all the functionalities
type ExecutionDB interface {
	Find(executionHash string) (*execution.Execution, error)
	Save(execution *execution.Execution) error
	OpenTransaction() (ExecutionTransaction, error)
	io.Closer
}

// ExecutionTransaction is the transaction handle.
type ExecutionTransaction interface {
	Find(executionHash string) (*execution.Execution, error)
	Save(execution *execution.Execution) error
	Commit() error
	Discard()
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

// Find the execution based on an executionHash, returns an error if not found
func (db *LevelDBExecutionDB) Find(executionHash string) (*execution.Execution, error) {
	return executionFind(db.db, executionHash)
}

// Save an instance of executable in the database
// Returns an error if anything from marshaling to database saving goes wrong
func (db *LevelDBExecutionDB) Save(execution *execution.Execution) error {
	return executionSave(db.db, execution)
}

// OpenTransaction opens an atomic DB transaction. Only one transaction can be
// opened at a time.
func (db *LevelDBExecutionDB) OpenTransaction() (ExecutionTransaction, error) {
	tx, err := db.db.OpenTransaction()
	if err != nil {
		return nil, err
	}
	return &LevelDBExecutionTransaction{tx: tx}, nil
}

// Close closes database.
func (db *LevelDBExecutionDB) Close() error {
	return db.db.Close()
}

// LevelDBExecutionTransaction is the transaction handle.
type LevelDBExecutionTransaction struct {
	tx *leveldb.Transaction
}

// Find the execution based on an executionHash, returns an error if not found
func (tx *LevelDBExecutionTransaction) Find(executionHash string) (*execution.Execution, error) {
	return executionFind(tx.tx, executionHash)
}

// Save an instance of executable in the database
// Returns an error if anything from marshaling to database saving goes wrong
func (tx *LevelDBExecutionTransaction) Save(execution *execution.Execution) error {
	return executionSave(tx.tx, execution)
}

// Commit commits the transaction.
func (tx *LevelDBExecutionTransaction) Commit() error {
	return tx.tx.Commit()
}

// Discard discards the transaction.
func (tx *LevelDBExecutionTransaction) Discard() {
	tx.tx.Discard()
}

// Find the execution based on an executionHash, returns an error if not found
func executionFind(db leveldbTxDB, executionHash string) (*execution.Execution, error) {
	data, err := db.Get([]byte(executionHash), nil)
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
func executionSave(db leveldbTxDB, execution *execution.Execution) error {
	if execution.Hash == "" {
		return errors.New("database: can't save service without id")
	}
	data, err := json.Marshal(execution)
	if err != nil {
		return err
	}
	return db.Put([]byte(execution.Hash), data, nil)
}
