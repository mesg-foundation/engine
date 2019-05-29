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
	e := &execution.Execution{Hash: []byte{'1'}}
	db.Save(e)
	tests := []struct {
		hash     []byte
		hasError bool
	}{
		{hash: e.Hash, hasError: false},
		{hash: []byte{1}, hasError: true},
	}
	for _, test := range tests {
		execution, err := db.Find(test.hash)
		if test.hasError {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		require.NotNil(t, execution)
		e, err := db.Find(execution.Hash)
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
		{&execution.Execution{Hash: []byte{'1'}}, false},
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
