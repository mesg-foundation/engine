package service

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

var serversubmit = new(Server)

func execute(name string) *execution.Execution {
	var inputs map[string]interface{}
	e, _ := execution.Create(&service.Service{
		Name: name,
		Tasks: map[string]*service.Task{
			"test": {
				Outputs: map[string]*service.Output{
					"output": {},
				},
			},
		},
	}, "test", inputs)
	e.Execute()
	return e
}

func TestSubmit(t *testing.T) {
	execution := execute("TestSubmit")
	reply, err := serversubmit.SubmitResult(context.Background(), &SubmitResultRequest{
		ExecutionID: execution.ID,
		OutputKey:   "output",
		OutputData:  "{}",
	})

	require.Nil(t, err)
	require.NotNil(t, reply)
}

func TestSubmitWithInvalidJSON(t *testing.T) {
	execution := execute("TestSubmitWithInvalidJSON")
	_, err := serversubmit.SubmitResult(context.Background(), &SubmitResultRequest{
		ExecutionID: execution.ID,
		OutputKey:   "output",
		OutputData:  "",
	})

	require.NotNil(t, err)
}

func TestSubmitWithInvalidID(t *testing.T) {
	_, err := serversubmit.SubmitResult(context.Background(), &SubmitResultRequest{
		ExecutionID: "xxxx",
		OutputKey:   "output",
		OutputData:  "",
	})
	require.NotNil(t, err)
	x, missingExecutionError := err.(*MissingExecutionError)
	require.True(t, missingExecutionError)
	require.Equal(t, "xxxx", x.ID)
}
