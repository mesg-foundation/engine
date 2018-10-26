package execution

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mesg-foundation/core/service"
)

var (
	task = &service.Task{
		Key: "task",
		Inputs: []*service.Parameter{
			&service.Parameter{Key: "foo", Type: "String"},
			&service.Parameter{Key: "bar", Type: "String"},
		},
	}
	defaultInputs = map[string]interface{}{
		"foo": "hello",
		"bar": "world",
	}
	tags = []string{"tag1", "tag2"}
)

func db(t *testing.T, dir string) DB {
	fmt.Println(dir)
	db, err := New(dir)
	require.NoError(t, err)
	return db
}

func TestCreate(t *testing.T) {
	db := db(t, filepath.Join(os.TempDir(), "TestCreate"))
	tests := []struct {
		inputs map[string]interface{}
		assert bool
	}{
		{inputs: map[string]interface{}{"foo": "hello", "bar": "world"}, assert: false},
		{inputs: map[string]interface{}{"foo": "hello"}, assert: true},
	}
	for _, test := range tests {
		execution, err := db.Create(task, test.inputs, tags)
		if test.assert {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.NotNil(t, execution)
			e, err := db.Find(execution.ID)
			require.NoError(t, err)
			require.NotNil(t, e)
			require.Equal(t, task, e.Task)
			require.Equal(t, test.inputs, e.Inputs)
			require.Equal(t, tags, e.Tags)
			require.Equal(t, execution.Status, Created)
			require.NotZero(t, e.CreatedAt)
		}
	}
}

func TestFind(t *testing.T) {
	db := db(t, filepath.Join(os.TempDir(), "TestFindExecution"))
	e, _ := db.Create(task, defaultInputs, tags)
	tests := []struct {
		id     []byte
		assert bool
	}{
		{e.ID, false},
		{[]byte("noid"), true},
	}
	for _, test := range tests {
		e, err := db.Find(test.id)
		if test.assert {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.NotNil(t, e)
		}
	}
}

func TestExecute(t *testing.T) {
	db := db(t, filepath.Join(os.TempDir(), "TestExecute"))
	e, _ := db.Create(&service.Task{Key: "TestExecute"}, map[string]interface{}{}, tags)
	tests := []struct {
		id     []byte
		assert bool
	}{
		{e.ID, false},
		{[]byte("doesn't exists"), true},
		{e.ID, true}, // this one is already executed so it should return an error
	}
	for _, test := range tests {
		e, err := db.Execute(test.id)
		if test.assert {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.NotNil(t, e)
			e, err := db.Find(e.ID)
			require.NoError(t, err)
			require.NotNil(t, e)
			require.Equal(t, e.Status, InProgress)
			require.NotZero(t, e.ExecutedAt)
		}
	}
}

func TestComplete(t *testing.T) {
	db := db(t, filepath.Join(os.TempDir(), "TestComplete"))
	e, _ := db.Create(&service.Task{Key: "TestComplete", Outputs: []*service.Output{
		&service.Output{
			Key: "outputX",
			Data: []*service.Parameter{
				&service.Parameter{
					Key:  "foo",
					Type: "String",
				},
			},
		},
	}}, map[string]interface{}{}, tags)
	db.Execute(e.ID)
	tests := []struct {
		id     []byte
		key    string
		data   map[string]interface{}
		assert bool
	}{
		{id: []byte("doesn't exists"), key: "", data: map[string]interface{}{}, assert: true},
		{id: e.ID, key: "output", data: map[string]interface{}{"foo": "bar"}, assert: true},
		{id: e.ID, key: "outputX", data: map[string]interface{}{}, assert: true},
		{id: e.ID, key: "outputX", data: map[string]interface{}{"foo": "bar"}, assert: false},
		{id: e.ID, key: "outputX", data: map[string]interface{}{"foo": "bar"}, assert: true}, // this one is already proccessed
	}
	for _, test := range tests {
		e, err := db.Complete(test.id, test.key, test.data)
		if test.assert {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.NotNil(t, e)
			e, err := db.Find(e.ID)
			require.NoError(t, err)
			require.NotNil(t, e)
			require.Equal(t, e.Status, Completed)
			require.NotZero(t, e.ExecutionDuration)
		}
	}
}

func TestConsistentID(t *testing.T) {
	db := db(t, filepath.Join(os.TempDir(), "TestConsistentID"))
	e, _ := db.Create(&service.Task{Key: "TestConsistentID", Outputs: []*service.Output{
		&service.Output{Key: "foo"},
	}}, map[string]interface{}{}, tags)
	require.NotZero(t, e.ID)
	e2, _ := db.Execute(e.ID)
	require.Equal(t, e.ID, e2.ID)
	e3, _ := db.Complete(e2.ID, "foo", map[string]interface{}{})
	require.Equal(t, e.ID, e3.ID)
}
