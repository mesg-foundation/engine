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
	e, _ := db.Save(&execution.Execution{})
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
		{&execution.Execution{}, false},
	}
	for _, test := range tests {
		e, err := db.Save(test.execution)
		if test.hasError {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		require.NotNil(t, e)
	}
}
