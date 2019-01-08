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
	"io/ioutil"
	"os"
	"testing"

	"github.com/mesg-foundation/core/execution"
	"github.com/stretchr/testify/require"
)

func db(t *testing.T, dir string) ExecutionDB {
	db, err := NewExecutionDB(dir)
	require.NoError(t, err)
	return db
}

func TestFind(t *testing.T) {
	dir, _ := ioutil.TempDir("", "TestFind")
	defer os.RemoveAll(dir)
	db := db(t, dir)
	defer db.Close()
	e := &execution.Execution{ID: "xxx"}
	db.Save(e)
	tests := []struct {
		id       string
		hasError bool
	}{
		{id: e.ID, hasError: false},
		{id: "doesn't exists", hasError: true},
	}
	for _, test := range tests {
		execution, err := db.Find(test.id)
		if test.hasError {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		require.NotNil(t, execution)
		e, err := db.Find(execution.ID)
		require.NoError(t, err)
		require.NotNil(t, e)
	}
}

func TestSave(t *testing.T) {
	dir, _ := ioutil.TempDir("", "TestSave")
	defer os.RemoveAll(dir)
	db := db(t, dir)
	defer db.Close()
	tests := []struct {
		execution *execution.Execution
		hasError  bool
	}{
		{&execution.Execution{ID: "xxx"}, false},
		{&execution.Execution{}, true},
	}
	for _, test := range tests {
		err := db.Save(test.execution)
		if test.hasError {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
	}
}
