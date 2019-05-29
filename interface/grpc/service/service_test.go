package service

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/protobuf/serviceapi"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestEmit(t *testing.T) {
	var (
		eventKey       = "request"
		eventData      = `{"data":{}}`
		server, closer = newServer(t)
	)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, eventServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Hash, false)

	ln, err := server.api.ListenEvent(s.Hash, nil)
	require.NoError(t, err)
	defer ln.Close()

	_, err = server.EmitEvent(context.Background(), &serviceapi.EmitEventRequest{
		Token:     s.Hash,
		EventKey:  eventKey,
		EventData: eventData,
	})
	require.NoError(t, err)

	event := <-ln.C
	require.Equal(t, eventKey, event.Key)
	require.Equal(t, eventData, jsonMarshal(t, event.Data))
}

func TestEmitNoData(t *testing.T) {
	var (
		eventKey       = "request"
		server, closer = newServer(t)
	)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, eventServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Hash, false)

	_, err = server.EmitEvent(context.Background(), &serviceapi.EmitEventRequest{
		Token:    s.Hash,
		EventKey: eventKey,
	})
	require.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestEmitWrongData(t *testing.T) {
	var (
		eventKey       = "request"
		server, closer = newServer(t)
	)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, eventServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Hash, false)

	_, err = server.EmitEvent(context.Background(), &serviceapi.EmitEventRequest{
		Token:     s.Hash,
		EventKey:  eventKey,
		EventData: "",
	})
	require.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestEmitWrongEvent(t *testing.T) {
	var (
		eventKey       = "test"
		server, closer = newServer(t)
	)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, eventServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Hash, false)

	_, err = server.EmitEvent(context.Background(), &serviceapi.EmitEventRequest{
		Token:     s.Hash,
		EventKey:  eventKey,
		EventData: "{}",
	})
	require.Error(t, err)
	notFoundErr, ok := err.(*service.EventNotFoundError)
	require.True(t, ok)
	require.Equal(t, eventKey, notFoundErr.EventKey)
	require.Equal(t, s.Name, notFoundErr.ServiceName)
}

func TestEmitInvalidData(t *testing.T) {
	var (
		eventKey       = "request"
		eventData      = `{"body":{}}`
		server, closer = newServer(t)
	)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, eventServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Hash, false)

	_, err = server.EmitEvent(context.Background(), &serviceapi.EmitEventRequest{
		Token:     s.Hash,
		EventKey:  eventKey,
		EventData: eventData,
	})
	require.Error(t, err)
	invalidErr, ok := err.(*service.InvalidEventDataError)
	require.True(t, ok)
	require.Equal(t, eventKey, invalidErr.EventKey)
	require.Equal(t, s.Name, invalidErr.ServiceName)
}

func TestServiceNotExists(t *testing.T) {
	server, closer := newServer(t)
	defer closer()

	_, err := server.EmitEvent(context.Background(), &serviceapi.EmitEventRequest{
		Token:     "TestServiceNotExists",
		EventKey:  "test",
		EventData: "{}",
	})
	require.Error(t, err)
}

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

	executionHash, err := server.api.ExecuteTask(s.Hash, taskKey, taskData, nil)
	require.NoError(t, err)

	ef := &api.ExecutionFilter{
		Statuses: []execution.Status{execution.Completed},
	}
	ln, err := server.api.ListenExecution(s.Hash, ef)
	require.NoError(t, err)
	defer ln.Close()

	_, err = server.SubmitResult(context.Background(), &serviceapi.SubmitResultRequest{
		ExecutionHash: string(executionHash),
		Result: &serviceapi.SubmitResultRequest_OutputData{
			OutputData: outputData,
		},
	})
	require.NoError(t, err)

	execution := <-ln.C
	require.Equal(t, executionHash, execution.Hash)
	require.Equal(t, outputData, jsonMarshal(t, execution.Outputs))
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

	executionHash, err := server.api.ExecuteTask(s.Hash, taskKey, taskData, nil)
	require.NoError(t, err)

	_, err = server.SubmitResult(context.Background(), &serviceapi.SubmitResultRequest{
		ExecutionHash: string(executionHash),
		Result:        &serviceapi.SubmitResultRequest_OutputData{},
	})
	require.Contains(t, err.Error(), "unexpected end of JSON input")
}

func TestSubmitWithInvalidID(t *testing.T) {
	var (
		outputData     = "{}"
		executionHash  = "1"
		server, closer = newServer(t)
	)
	defer closer()

	_, err := server.SubmitResult(context.Background(), &serviceapi.SubmitResultRequest{
		ExecutionHash: executionHash,
		Result: &serviceapi.SubmitResultRequest_OutputData{
			OutputData: outputData,
		},
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

	executionHash, err := server.api.ExecuteTask(s.Hash, taskKey, taskData, nil)
	require.NoError(t, err)

	_, err = server.SubmitResult(context.Background(), &serviceapi.SubmitResultRequest{
		ExecutionHash: string(executionHash),
		Result: &serviceapi.SubmitResultRequest_OutputData{
			OutputData: outputData,
		},
	})
	require.Error(t, err)
	invalidErr, ok := err.(*service.InvalidTaskOutputError)
	require.True(t, ok)
	require.Equal(t, taskKey, invalidErr.TaskKey)
	require.Equal(t, s.Name, invalidErr.ServiceName)
}
