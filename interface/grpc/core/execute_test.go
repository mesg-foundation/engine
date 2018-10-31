package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	var (
		taskKey        = "call"
		data           = `{"url": "https://mesg.tech", "data": {}, "headers": {}}`
		server, closer = newServer(t)
	)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, taskServicePath))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.ID)

	require.NoError(t, server.api.StartService(s.ID))
	defer server.api.StopService(s.ID)

	reply, err := server.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID: s.ID,
		TaskKey:   taskKey,
		InputData: data,
	})
	require.NoError(t, err)
	require.NotEqual(t, "", reply.ExecutionID)
}

func TestExecuteWithInvalidJSON(t *testing.T) {
	server, closer := newServer(t)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, taskServicePath))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.ID)

	_, err = server.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID: s.ID,
		TaskKey:   "test",
		InputData: "",
	})
	require.Error(t, err)
	require.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestExecuteWithInvalidTask(t *testing.T) {
	var (
		taskKey        = "error"
		server, closer = newServer(t)
	)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, taskServicePath))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.ID)

	require.NoError(t, server.api.StartService(s.ID))
	defer server.api.StopService(s.ID)

	_, err = server.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID: s.ID,
		TaskKey:   taskKey,
		InputData: "{}",
	})
	require.Error(t, err)
	notFoundErr, ok := err.(*service.TaskNotFoundError)
	require.True(t, ok)
	require.Equal(t, taskKey, notFoundErr.TaskKey)
	require.Equal(t, s.Name, notFoundErr.ServiceName)
}

func TestExecuteWithInvalidTaskInput(t *testing.T) {
	var (
		taskKey        = "call"
		data           = `{"headers": {}}`
		server, closer = newServer(t)
	)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, taskServicePath))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.ID)

	require.NoError(t, server.api.StartService(s.ID))
	defer server.api.StopService(s.ID)

	_, err = server.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID: s.ID,
		TaskKey:   taskKey,
		InputData: data,
	})
	require.Error(t, err)
	invalidErr, ok := err.(*service.InvalidTaskInputError)
	require.True(t, ok)
	require.Equal(t, taskKey, invalidErr.TaskKey)
	require.Equal(t, s.Name, invalidErr.ServiceName)
}

func TestExecuteWithNonRunningService(t *testing.T) {
	server, closer := newServer(t)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, taskServicePath))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.ID)

	_, err = server.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID: s.ID,
		TaskKey:   "test",
		InputData: "{}",
	})
	require.Equal(t, &api.NotRunningServiceError{ServiceID: s.ID}, err)
}

func TestExecuteWithNonExistingService(t *testing.T) {
	server, closer := newServer(t)
	defer closer()

	_, err := server.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID: "service that doesnt exists",
		TaskKey:   "error",
		InputData: "{}",
	})
	require.Error(t, err)
}
