package client

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

var serverexecute = new(Server)

var testService = service.Service{
	Name: "TestService",
	Tasks: map[string]*service.Task{
		"test": &service.Task{},
	},
}

func TestExecute(t *testing.T) {
	reply, err := serverexecute.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		Service:  &testService,
		TaskKey:  "test",
		TaskData: "{}",
	})

	assert.Nil(t, err)
	assert.NotNil(t, reply)
}

func TestExecuteWithInvalidJSON(t *testing.T) {
	_, err := serverexecute.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		Service:  &testService,
		TaskKey:  "test",
		TaskData: "",
	})

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestExecuteWithInvalidTask(t *testing.T) {
	_, err := serverexecute.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		Service:  &testService,
		TaskKey:  "error",
		TaskData: "{}",
	})

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Task error doesn't exists in service TestService")
}
