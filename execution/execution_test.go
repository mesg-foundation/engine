package execution

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mesg-foundation/core/service"
)

var (
	srv = &service.Service{
		Tasks: []*service.Task{
			&service.Task{
				Key: "task",
				Inputs: []*service.Parameter{
					&service.Parameter{Key: "foo", Type: "String"},
					&service.Parameter{Key: "bar", Type: "String"},
				},
				Outputs: []*service.Output{
					&service.Output{
						Key: "outputX",
						Data: []*service.Parameter{
							&service.Parameter{
								Key:  "foo",
								Type: "String",
							},
						},
					},
				},
			},
		},
	}
	taskKey       = "task"
	defaultInputs = map[string]interface{}{
		"foo": "hello",
		"bar": "world",
	}
	tags = []string{"tag1", "tag2"}
)

func db(t *testing.T, dir string) DB {
	db, err := New(dir)
	require.NoError(t, err)
	return db
}

func TestCreate(t *testing.T) {
	dir, _ := ioutil.TempDir("", "TestCreate")
	defer os.RemoveAll(dir)
	db := db(t, dir)
	defer db.Close()
	tests := []struct {
		inputs map[string]interface{}
		assert bool
	}{
		{inputs: map[string]interface{}{"foo": "hello", "bar": "world"}, assert: false},
		{inputs: map[string]interface{}{"foo": "hello"}, assert: true},
	}
	for _, test := range tests {
		execution, err := db.Create(srv, taskKey, test.inputs, tags)
		if test.assert {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.NotNil(t, execution)
			e, err := db.Find(execution.ID)
			require.NoError(t, err)
			require.NotNil(t, e)
			require.Equal(t, srv.ID, e.Service.ID)
			require.Equal(t, taskKey, e.TaskKey)
			require.Equal(t, test.inputs, e.Inputs)
			require.Equal(t, tags, e.Tags)
			require.Equal(t, execution.Status, Created)
			require.NotZero(t, e.CreatedAt)
		}
	}
}

func TestFind(t *testing.T) {
	dir, _ := ioutil.TempDir("", "TestFindExecution")
	defer os.RemoveAll(dir)
	db := db(t, dir)
	defer db.Close()
	e, _ := db.Create(srv, taskKey, defaultInputs, tags)
	tests := []struct {
		id     string
		assert bool
	}{
		{e.ID, false},
		{"noid", true},
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
	dir, _ := ioutil.TempDir("", "TestExecute")
	defer os.RemoveAll(dir)
	db := db(t, dir)
	defer db.Close()
	e, _ := db.Create(srv, taskKey, map[string]interface{}{"foo": "1", "bar": "2"}, tags)
	tests := []struct {
		id     string
		assert bool
	}{
		{e.ID, false},
		{"doesn't exists", true},
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
	dir, _ := ioutil.TempDir("", "TestComplete")
	defer os.RemoveAll(dir)
	db := db(t, dir)
	defer db.Close()
	e, _ := db.Create(srv, taskKey, map[string]interface{}{"foo": "1", "bar": "2"}, tags)
	db.Execute(e.ID)
	tests := []struct {
		id     string
		key    string
		data   map[string]interface{}
		assert bool
	}{
		{id: "doesn't exists", key: "", data: map[string]interface{}{}, assert: true},
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
	dir, _ := ioutil.TempDir("", "TestConsistentID")
	defer os.RemoveAll(dir)
	db := db(t, dir)
	defer db.Close()
	e, _ := db.Create(srv, taskKey, map[string]interface{}{"foo": "1", "bar": "2"}, tags)
	require.NotZero(t, e.ID)
	e2, _ := db.Execute(e.ID)
	require.Equal(t, e.ID, e2.ID)
	e3, _ := db.Complete(e2.ID, "outputX", map[string]interface{}{"foo": "1", "bar": "2"})
	require.Equal(t, e.ID, e3.ID)
}
