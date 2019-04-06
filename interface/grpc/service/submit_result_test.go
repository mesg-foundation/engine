package service

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/protobuf/serviceapi"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestSubmit(t *testing.T) {
	var (
		taskKey  = "call"
		taskData = map[string]interface{}{
			"url":     "https://mesg.com",
			"data":    map[string]interface{}{},
			"headers": map[string]interface{}{},
		}
		outputData     = `{"foo":{}}`
		server, closer = newServer(t)
	)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, taskServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Hash, false)

	require.NoError(t, server.api.StartService(s.Hash))
	defer server.api.StopService(s.Hash)

	executionID, err := server.api.ExecuteTask(s.Hash, taskKey, taskData, nil)
	require.NoError(t, err)

	ln, err := server.api.ListenResult(s.Hash)
	require.NoError(t, err)
	defer ln.Close()

	_, err = server.SubmitResult(context.Background(), &serviceapi.SubmitResultRequest{
		ExecutionID: executionID,
		Data:        outputData,
	})
	require.NoError(t, err)

	select {
	case err := <-ln.Err:
		t.Error(err)

	case execution := <-ln.Executions:
		require.Equal(t, executionID, execution.ID)
		require.Equal(t, outputData, jsonMarshal(t, execution.Outputs))
	}
}

func TestSubmitWithInvalidJSON(t *testing.T) {
	var (
		taskKey  = "call"
		taskData = map[string]interface{}{
			"url":     "https://mesg.com",
			"data":    map[string]interface{}{},
			"headers": map[string]interface{}{},
		}
		server, closer = newServer(t)
	)
	defer closer()
	s, validationErr, err := server.api.DeployService(serviceTar(t, taskServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Hash, false)

	require.NoError(t, server.api.StartService(s.Hash))
	defer server.api.StopService(s.Hash)

	executionID, err := server.api.ExecuteTask(s.Hash, taskKey, taskData, nil)
	require.NoError(t, err)

	_, err = server.SubmitResult(context.Background(), &serviceapi.SubmitResultRequest{
		ExecutionID: executionID,
		Data:        "",
	})
	require.Equal(t, "invalid output data error: unexpected end of JSON input", err.Error())
}

func TestSubmitWithInvalidID(t *testing.T) {
	var (
		outputData     = "{}"
		executionID    = "1"
		server, closer = newServer(t)
	)
	defer closer()

	_, err := server.SubmitResult(context.Background(), &serviceapi.SubmitResultRequest{
		ExecutionID: executionID,
		Data:        outputData,
	})
	require.Error(t, err)
}

func TestSubmitWithInvalidTaskOutputs(t *testing.T) {
	var (
		taskKey  = "call"
		taskData = map[string]interface{}{
			"url":     "https://mesg.com",
			"data":    map[string]interface{}{},
			"headers": map[string]interface{}{},
		}
		outputData     = `{"foo":1}`
		server, closer = newServer(t)
	)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, taskServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Hash, false)

	require.NoError(t, server.api.StartService(s.Hash))
	defer server.api.StopService(s.Hash)

	executionID, err := server.api.ExecuteTask(s.Hash, taskKey, taskData, nil)
	require.NoError(t, err)

	_, err = server.SubmitResult(context.Background(), &serviceapi.SubmitResultRequest{
		ExecutionID: executionID,
		Data:        outputData,
	})
	require.Error(t, err)
	invalidErr, ok := err.(*service.InvalidTaskOutputError)
	require.True(t, ok)
	require.Equal(t, taskKey, invalidErr.TaskKey)
	require.Equal(t, s.Name, invalidErr.ServiceName)
}
