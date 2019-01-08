// Copyright 2018 MESG Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	data, err := db.db.Get([]byte(executionID), nil)
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
	if execution.ID == "" {
		return errors.New("database: can't save service without id")
	}
	data, err := json.Marshal(execution)
	if err != nil {
		return err
	}
	return db.db.Put([]byte(execution.ID), data, nil)
}

// Close closes database.
func (db *LevelDBExecutionDB) Close() error {
	return db.db.Close()
}
