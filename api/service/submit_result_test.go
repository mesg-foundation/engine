package service

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/service"

	"github.com/stvp/assert"
)

var serversubmit = new(Server)

func execute(name string) (e *execution.Execution) {
	var inputs interface{}
	e, _ = execution.Create(&service.Service{
		Name: name,
		Tasks: map[string]*service.Task{
			"test": &service.Task{},
		},
	}, "test", inputs)
	e.Execute()
	return
}

func TestSubmit(t *testing.T) {
	execution := execute("TestSubmit")
	reply, err := serversubmit.SubmitResult(context.Background(), &SubmitResultRequest{
		ExecutionID: execution.ID,
		OutputKey:   "output",
		OutputData:  "{}",
	})

	assert.Nil(t, err)
	assert.NotNil(t, reply)
}

func TestSubmitWithInvalidJSON(t *testing.T) {
	execution := execute("TestSubmitWithInvalidJSON")
	_, err := serversubmit.SubmitResult(context.Background(), &SubmitResultRequest{
		ExecutionID: execution.ID,
		OutputKey:   "output",
		OutputData:  "",
	})

	assert.NotNil(t, err)
}

func TestSubmitWithInvalidID(t *testing.T) {
	_, err := serversubmit.SubmitResult(context.Background(), &SubmitResultRequest{
		ExecutionID: "xxxx",
		OutputKey:   "output",
		OutputData:  "",
	})

	assert.NotNil(t, err)
}
