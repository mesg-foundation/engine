package result

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/service"

	"github.com/mesg-foundation/core/types"
	"github.com/stvp/assert"
)

var serversubmit = new(Server)

func execute(name string) (reply *types.TaskReply) {
	var inputs interface{}
	execution, _ := execution.Create(&service.Service{
		Name: name,
		Tasks: map[string]*types.ProtoTask{
			"test": &types.ProtoTask{},
		},
	}, "test", inputs)
	reply, _ = execution.Execute()
	return
}

func TestSubmit(t *testing.T) {
	execution := execute("TestSubmit")
	reply, err := serversubmit.Submit(context.Background(), &types.SubmitResultRequest{
		ExecutionID: execution.ExecutionID,
		Output:      "output",
		Task:        execution.Task,
		Data:        "{}",
	})

	assert.Nil(t, err)
	assert.NotNil(t, reply)
}

func TestSubmitWithInvalidJSON(t *testing.T) {
	execution := execute("TestSubmitWithInvalidJSON")
	_, err := serversubmit.Submit(context.Background(), &types.SubmitResultRequest{
		ExecutionID: execution.ExecutionID,
		Output:      "output",
		Task:        execution.Task,
		Data:        "",
	})

	assert.NotNil(t, err)
}

func TestSubmitWithInvalidID(t *testing.T) {
	execution := execute("TestSubmitWithInvalidID")
	_, err := serversubmit.Submit(context.Background(), &types.SubmitResultRequest{
		ExecutionID: "xxxx",
		Output:      "output",
		Task:        execution.Task,
		Data:        "",
	})

	assert.NotNil(t, err)
}
