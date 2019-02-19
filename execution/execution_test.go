package execution

import (
	"encoding/hex"
	"errors"
	"testing"
	"testing/quick"

	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/utils/structhash"
	"github.com/stretchr/testify/require"
)

var (
	serviceName      = "SERVICE_ID"
	eventID          = "EVENT_ID"
	taskKey          = "TASK_KEY"
	taskKeyWithError = "TASK_KEY_ERROR"
	tags             = []string{"tag1", "tag2"}
	srv, _           = service.FromService(&service.Service{
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
						Key: "OUTPUT_KEY_1",
						Data: []*service.Parameter{
							{
								Key:  "foo",
								Type: "String",
							},
						},
					},
				},
			},
			{
				Key: taskKeyWithError,
				Inputs: []*service.Parameter{
					{Key: "foo", Type: "String"},
				},
				Outputs: []*service.Output{
					{
						Key: "OUTPUT_KEY_1",
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
		{
			name:    "success",
			taskKey: taskKey,
			inputs:  map[string]interface{}{"foo": "hello", "bar": "world"},
			err:     nil,
		},
		{
			name:    "task not found",
			taskKey: "wrongtask",
			inputs:  map[string]interface{}{},
			err: &service.TaskNotFoundError{
				TaskKey:     "wrongtask",
				ServiceName: serviceName,
			},
		},
		{
			name:    "invalid task input",
			taskKey: taskKey,
			inputs:  map[string]interface{}{"foo": "hello"},
			err: &service.InvalidTaskInputError{
				TaskKey:     taskKey,
				ServiceName: serviceName,
				Warnings: []*service.ParameterWarning{
					{
						Key:       "bar",
						Warning:   "required",
						Parameter: &service.Parameter{Key: "bar", Type: "String"},
					},
				},
			},
		},
	}
	for _, test := range tests {
		execution, err := New(srv, eventID, test.taskKey, test.inputs, tags)
		require.Equal(t, test.err, err, test.name)
		if test.err != nil {
			continue
		}
		require.NotNil(t, execution, test.name)
		require.Equal(t, srv.Hash, execution.Service.Hash, test.name)
		require.Equal(t, eventID, execution.EventID, test.name)
		require.Equal(t, taskKey, execution.TaskKey, test.name)
		require.Equal(t, test.inputs, execution.Inputs, test.name)
		require.Equal(t, tags, execution.Tags, test.name)
		require.Equal(t, Created, execution.Status, test.name)
		require.NotZero(t, execution.CreatedAt, test.name)
	}
}

func TestExecute(t *testing.T) {
	e, _ := New(srv, eventID, taskKey, map[string]interface{}{"foo": "1", "bar": "2"}, tags)
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "success",
			err:  nil,
		},
		{ // this one is already executed so it should return an error
			name: "status error",
			err: StatusError{
				ExpectedStatus: Created,
				ActualStatus:   InProgress,
			},
		},
	}
	for _, test := range tests {
		err := e.Execute()
		require.Equal(t, test.err, err, test.name)
		if test.err != nil {
			continue
		}
		require.NotNil(t, e, test.name)
		require.Equal(t, InProgress, e.Status, test.name)
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
		{
			name: "task output not found because of empty output key",
			key:  "",
			data: map[string]interface{}{},
			err: &service.TaskOutputNotFoundError{
				TaskKey:       taskKey,
				TaskOutputKey: "",
				ServiceName:   serviceName,
			},
		},
		{
			name: "task output not found because wrong output key",
			key:  "output",
			data: map[string]interface{}{"foo": "bar"},
			err: &service.TaskOutputNotFoundError{
				TaskKey:       taskKey,
				TaskOutputKey: "output",
				ServiceName:   serviceName,
			},
		},
		{
			name: "invalid task output",
			key:  "OUTPUT_KEY_1",
			data: map[string]interface{}{},
			err: &service.InvalidTaskOutputError{
				TaskKey:       taskKey,
				TaskOutputKey: "OUTPUT_KEY_1",
				ServiceName:   serviceName,
				Warnings: []*service.ParameterWarning{
					{
						Key:       "foo",
						Warning:   "required",
						Parameter: &service.Parameter{Key: "foo", Type: "String"},
					},
				},
			},
		},
		{
			name: "success",
			key:  "OUTPUT_KEY_1",
			data: map[string]interface{}{"foo": "bar"},
			err:  nil,
		},
		{ // this one is already proccessed
			name: "already executed",
			key:  "OUTPUT_KEY_1",
			data: map[string]interface{}{"foo": "bar"},
			err: StatusError{
				ExpectedStatus: InProgress,
				ActualStatus:   Completed,
			},
		},
	}
	for _, test := range tests {
		err := e.Complete(test.key, test.data)
		require.Equal(t, test.err, err, test.name)
		if test.err != nil {
			continue
		}
		require.Equal(t, Completed, e.Status, test.name)
		require.NotZero(t, e.ExecutionDuration, test.name)
		require.Zero(t, e.Error, test.name)
	}
}

func TestFailed(t *testing.T) {
	e, _ := New(srv, eventID, taskKeyWithError, map[string]interface{}{"foo": "1", "bar": "2"}, tags)
	e.Execute()
	tests := []struct {
		name string
		xerr error
		err  error
	}{
		{
			name: "with failed error",
			xerr: errors.New("failed"),
		},
		{ // this one is already proccessed
			name: "with status error",
			err: StatusError{
				ExpectedStatus: InProgress,
				ActualStatus:   Failed,
			},
		},
	}
	for _, test := range tests {
		err := e.Failed(test.xerr)
		require.Equal(t, test.err, err, test.name)
		if test.err != nil {
			continue
		}
		require.Equal(t, Failed, e.Status, test.name)
		require.NotZero(t, e.ExecutionDuration, test.name)
		require.Equal(t, test.xerr.Error(), e.Error, test.name)
	}
}

func TestStatus(t *testing.T) {
	require.Equal(t, "created", Created.String())
	require.Equal(t, "in progress", InProgress.String())
	require.Equal(t, "completed", Completed.String())
	require.Equal(t, "failed", Failed.String())
}

func TestExecutionIDHash(t *testing.T) {
	ids := make(map[string]bool)

	f := func(eventID, taskKey, inputs string, tags []string) bool {
		e := Execution{
			EventID: eventID,
			TaskKey: taskKey,
			Tags:    tags,
			Inputs:  map[string]interface{}{inputs: "input"},
		}

		h := structhash.Sha1(e)
		id := hex.EncodeToString(h[:])
		if ids[id] {
			return false
		}
		ids[id] = true
		return true
	}

	require.NoError(t, quick.Check(f, nil))
}
