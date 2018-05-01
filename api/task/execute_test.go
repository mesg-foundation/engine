package task

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/types"
	"github.com/stvp/assert"
)

var serverexecute = new(Server)

var testService = types.ProtoService{
	Name: "TestService",
	Tasks: map[string]*types.ProtoTask{
		"test": &types.ProtoTask{},
	},
}

func TestExecute(t *testing.T) {
	reply, err := serverexecute.Execute(context.Background(), &types.ExecuteTaskRequest{
		Service: &testService,
		Task:    "test",
		Data:    "{}",
	})

	assert.Nil(t, err)
	assert.NotNil(t, reply)
}

func TestExecuteWithInvalidJSON(t *testing.T) {
	_, err := serverexecute.Execute(context.Background(), &types.ExecuteTaskRequest{
		Service: &testService,
		Task:    "test",
		Data:    "",
	})

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestExecuteWithInvalidTask(t *testing.T) {
	_, err := serverexecute.Execute(context.Background(), &types.ExecuteTaskRequest{
		Service: &testService,
		Task:    "error",
		Data:    "{}",
	})

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Task error doesn't exists in service TestService")
}
