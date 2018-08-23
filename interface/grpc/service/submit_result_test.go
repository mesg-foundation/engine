package service

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

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
	}, "test", inputs, []string{})
	e.Execute()
	return e
}

func TestSubmit(t *testing.T) {
	server := newServer(t)
	execution := execute("TestSubmit")
	reply, err := server.SubmitResult(context.Background(), &SubmitResultRequest{
		ExecutionID: execution.ID,
		OutputKey:   "output",
		OutputData:  "{}",
	})

	require.Nil(t, err)
	require.NotNil(t, reply)
}

func TestSubmitWithInvalidJSON(t *testing.T) {
	server := newServer(t)
	execution := execute("TestSubmitWithInvalidJSON")
	_, err := server.SubmitResult(context.Background(), &SubmitResultRequest{
		ExecutionID: execution.ID,
		OutputKey:   "output",
		OutputData:  "",
	})

	require.NotNil(t, err)
}

func TestSubmitWithInvalidID(t *testing.T) {
	server := newServer(t)
	executionID := "xxxx"
	_, err := server.SubmitResult(context.Background(), &SubmitResultRequest{
		ExecutionID: executionID,
		OutputKey:   "output",
		OutputData:  "{}",
	})
	require.Equal(t, &api.MissingExecutionError{ID: executionID}, err)
}
