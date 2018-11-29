package execution

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mesg-foundation/core/service"
)

var (
	serviceName = "1"
	eventID     = "2"
	taskKey     = "task"
	tags        = []string{"tag1", "tag2"}
	srv, _      = service.FromService(&service.Service{
		Name: serviceName,
		Tasks: []*service.Task{
			{
				Key: taskKey,
				Inputs: []*service.Parameter{
					{Key: "foo", Type: "String"},
					{Key: "bar", Type: "String"},
				},
				Outputs: []*service.Output{
					{
						Key: "outputX",
						Data: []*service.Parameter{
							{
								Key:  "foo",
								Type: "String",
							},
						},
					},
				},
			},
		},
	})
)

func TestNewFromService(t *testing.T) {
	tests := []struct {
		name    string
		taskKey string
		inputs  map[string]interface{}
		err     error
	}{
		{name: "1", taskKey: taskKey, inputs: map[string]interface{}{"foo": "hello", "bar": "world"}, err: nil},
		{name: "2", taskKey: "wrongtask", inputs: map[string]interface{}{}, err: &service.TaskNotFoundError{
			TaskKey:     "wrongtask",
			ServiceName: serviceName,
		}},
		{name: "3", taskKey: taskKey, inputs: map[string]interface{}{"foo": "hello"}, err: &service.InvalidTaskInputError{
			TaskKey:     taskKey,
			ServiceName: serviceName,
			Warnings: []*service.ParameterWarning{
				{
					Key:       "bar",
					Warning:   "required",
					Parameter: &service.Parameter{Key: "bar", Type: "String"},
				},
			},
		}},
	}
	for _, test := range tests {
		execution, err := New(srv, eventID, test.taskKey, test.inputs, tags)
		require.Equal(t, test.err, err, test.name)
		if test.err != nil {
			continue
		}
		require.NotNil(t, execution, test.name)
		require.Equal(t, srv.ID, execution.Service.ID, test.name)
		require.Equal(t, eventID, execution.EventID, test.name)
		require.Equal(t, taskKey, execution.TaskKey, test.name)
		require.Equal(t, test.inputs, execution.Inputs, test.name)
		require.Equal(t, tags, execution.Tags, test.name)
		require.Equal(t, execution.Status, Created, test.name)
		require.NotZero(t, execution.CreatedAt, test.name)
	}
}

func TestExecute(t *testing.T) {
	e, _ := New(srv, eventID, taskKey, map[string]interface{}{"foo": "1", "bar": "2"}, tags)
	tests := []struct {
		name string
		err  error
	}{
		{name: "1", err: nil},
		{name: "2", err: StatusError{ExpectedStatus: Created, ActualStatus: InProgress}}, // this one is already executed so it should return an error
	}
	for _, test := range tests {
		err := e.Execute()
		require.Equal(t, test.err, err, test.name)
		if test.err != nil {
			continue
		}
		require.NotNil(t, e, test.name)
		require.Equal(t, e.Status, InProgress, test.name)
		require.NotZero(t, e.ExecutedAt, test.name)
	}
}

func TestComplete(t *testing.T) {
	e, _ := New(srv, eventID, taskKey, map[string]interface{}{"foo": "1", "bar": "2"}, tags)
	e.Execute()
	tests := []struct {
		name string
		key  string
		data map[string]interface{}
		err  error
	}{
		{name: "1", key: "", data: map[string]interface{}{}, err: &service.TaskOutputNotFoundError{
			TaskKey:       taskKey,
			TaskOutputKey: "",
			ServiceName:   serviceName},
		},
		{name: "2", key: "output", data: map[string]interface{}{"foo": "bar"}, err: &service.TaskOutputNotFoundError{
			TaskKey:       taskKey,
			TaskOutputKey: "output",
			ServiceName:   serviceName,
		}},
		{name: "3", key: "outputX", data: map[string]interface{}{}, err: &service.InvalidTaskOutputError{
			TaskKey:       taskKey,
			TaskOutputKey: "outputX",
			ServiceName:   serviceName,
			Warnings: []*service.ParameterWarning{
				{
					Key:       "foo",
					Warning:   "required",
					Parameter: &service.Parameter{Key: "foo", Type: "String"},
				},
			},
		}},
		{name: "4", key: "outputX", data: map[string]interface{}{"foo": "bar"}, err: nil},
		{name: "5", key: "outputX", data: map[string]interface{}{"foo": "bar"}, err: StatusError{
			ExpectedStatus: InProgress,
			ActualStatus:   Completed,
		}}, // this one is already proccessed
	}
	for _, test := range tests {
		err := e.Complete(test.key, test.data)
		require.Equal(t, test.err, err, test.name)
		if test.err != nil {
			continue
		}
		require.NotNil(t, e)
		require.NoError(t, err, test.name)
		require.NotNil(t, e, test.name)
		require.Equal(t, e.Status, Completed, test.name)
		require.NotZero(t, e.ExecutionDuration, test.name)
	}
}
