package execution

import (
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

func TestNewFromService(t *testing.T) {
	tests := []struct {
		taskKey  string
		inputs   map[string]interface{}
		hasError bool
	}{
		{taskKey: taskKey, inputs: map[string]interface{}{"foo": "hello", "bar": "world"}, hasError: false},
		{taskKey: "wrongtask", inputs: map[string]interface{}{}, hasError: true},
		{taskKey: taskKey, inputs: map[string]interface{}{"foo": "hello"}, hasError: true},
	}
	for _, test := range tests {
		execution, err := New(srv, test.taskKey, test.inputs, tags)
		if test.hasError {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		require.NotNil(t, execution)
		require.Equal(t, srv.ID, execution.Service.ID)
		require.Equal(t, taskKey, execution.TaskKey)
		require.Equal(t, test.inputs, execution.Inputs)
		require.Equal(t, tags, execution.Tags)
		require.Equal(t, execution.Status, Created)
		require.NotZero(t, execution.CreatedAt)
	}
}

func TestExecute(t *testing.T) {
	e, _ := New(srv, taskKey, map[string]interface{}{"foo": "1", "bar": "2"}, tags)
	tests := []struct {
		id       string
		hasError bool
	}{
		{e.ID, false},
		{"doesn't exists", true},
		{e.ID, true}, // this one is already executed so it should return an error
	}
	for _, test := range tests {
		err := e.Execute()
		if test.hasError {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		require.NotNil(t, e)
		require.NoError(t, err)
		require.NotNil(t, e)
		require.Equal(t, e.Status, InProgress)
		require.NotZero(t, e.ExecutedAt)
	}
}

func TestComplete(t *testing.T) {
	e, _ := New(srv, taskKey, map[string]interface{}{"foo": "1", "bar": "2"}, tags)
	e.Execute()
	tests := []struct {
		id       string
		key      string
		data     map[string]interface{}
		hasError bool
	}{
		{id: "doesn't exists", key: "", data: map[string]interface{}{}, hasError: true},
		{id: e.ID, key: "output", data: map[string]interface{}{"foo": "bar"}, hasError: true},
		{id: e.ID, key: "outputX", data: map[string]interface{}{}, hasError: true},
		{id: e.ID, key: "outputX", data: map[string]interface{}{"foo": "bar"}, hasError: false},
		{id: e.ID, key: "outputX", data: map[string]interface{}{"foo": "bar"}, hasError: true}, // this one is already proccessed
	}
	for _, test := range tests {
		err := e.Complete(test.key, test.data)
		if test.hasError {
			require.Error(t, err)
			continue
		}
		require.NoError(t, err)
		require.NotNil(t, e)
		require.NoError(t, err)
		require.NotNil(t, e)
		require.Equal(t, e.Status, Completed)
		require.NotZero(t, e.ExecutionDuration)
	}
}
