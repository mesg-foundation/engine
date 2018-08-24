package service

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/api"
	"github.com/stretchr/testify/require"
)

func TestSubmitA(t *testing.T) {
	var (
		path     = "./service-test-task"
		taskKey  = "call"
		taskData = map[string]interface{}{
			"url":     "https://mesg.tech",
			"data":    map[string]interface{}{},
			"headers": map[string]interface{}{},
		}
		outputKey  = "result"
		outputData = `{"data1":{}}`
		server     = newServer(t)
	)

	s, validationErr, err := server.api.DeployService(serviceTar(t, path))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Id)

	require.NoError(t, server.api.StartService(s.Id))
	defer server.api.StopService(s.Id)

	executionID, err := server.api.ExecuteTask(s.Id, taskKey, taskData, nil)
	require.NoError(t, err)

	ln, err := server.api.ListenResult(s.Id)
	require.NoError(t, err)
	defer ln.Close()

	_, err = server.SubmitResult(context.Background(), &SubmitResultRequest{
		ExecutionID: executionID,
		OutputKey:   outputKey,
		OutputData:  outputData,
	})
	require.NoError(t, err)

	select {
	case err := <-ln.Err:
		t.Error(err)

	case execution := <-ln.Executions:
		require.Equal(t, executionID, execution.ID)
		require.Equal(t, outputKey, execution.Output)
		require.Equal(t, outputData, jsonMarshal(t, execution.OutputData))
	}
}

func TestSubmitWithInvalidJSON(t *testing.T) {
	var (
		path     = "./service-test-task"
		taskKey  = "call"
		taskData = map[string]interface{}{
			"url":     "https://mesg.tech",
			"data":    map[string]interface{}{},
			"headers": map[string]interface{}{},
		}
		outputKey = "result"
		server    = newServer(t)
	)

	s, validationErr, err := server.api.DeployService(serviceTar(t, path))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Id)

	require.NoError(t, server.api.StartService(s.Id))
	defer server.api.StopService(s.Id)

	executionID, err := server.api.ExecuteTask(s.Id, taskKey, taskData, nil)
	require.NoError(t, err)

	_, err = server.SubmitResult(context.Background(), &SubmitResultRequest{
		ExecutionID: executionID,
		OutputKey:   outputKey,
		OutputData:  "",
	})
	require.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestSubmitWithInvalidID(t *testing.T) {
	var (
		outputKey   = "output"
		outputData  = "{}"
		executionID = "1"
		server      = newServer(t)
	)

	_, err := server.SubmitResult(context.Background(), &SubmitResultRequest{
		ExecutionID: executionID,
		OutputKey:   outputKey,
		OutputData:  outputData,
	})
	require.Equal(t, &api.MissingExecutionError{ID: executionID}, err)
}
