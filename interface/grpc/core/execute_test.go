package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	var (
		path    = "../../../service-test/task"
		taskKey = "call"
		data    = `{"url": "https://mesg.tech", "data": {}, "headers": {}}`
		server  = newServer(t)
	)

	s, validationErr, err := server.api.DeployService(serviceTar(t, path))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Id)

	require.NoError(t, server.api.StartService(s.Id))
	defer server.api.StopService(s.Id)

	reply, err := server.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		ServiceID: s.Id,
		TaskKey:   taskKey,
		InputData: data,
	})
	require.Nil(t, err)
	require.NotEqual(t, "", reply.ExecutionID)
}

func TestExecuteWithInvalidJSON(t *testing.T) {
	var (
		path   = "../../../service-test/task"
		server = newServer(t)
	)

	s, validationErr, err := server.api.DeployService(serviceTar(t, path))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Id)

	_, err = server.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		ServiceID: s.Id,
		TaskKey:   "test",
		InputData: "",
	})
	require.NotNil(t, err)
	require.Equal(t, err.Error(), "unexpected end of JSON input")
}

func TestExecuteWithInvalidTask(t *testing.T) {
	var (
		path    = "../../../service-test/task"
		taskKey = "error"
		server  = newServer(t)
	)

	s, validationErr, err := server.api.DeployService(serviceTar(t, path))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Id)

	require.NoError(t, server.api.StartService(s.Id))
	defer server.api.StopService(s.Id)

	_, err = server.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		ServiceID: s.Id,
		TaskKey:   taskKey,
		InputData: "{}",
	})
	require.Error(t, err)
	notFoundErr := err.(*service.TaskNotFoundError)
	require.Equal(t, taskKey, notFoundErr.TaskKey)
	require.Equal(t, s.Name, notFoundErr.ServiceName)
}

func TestExecuteWithInvalidTaskInput(t *testing.T) {
	var (
		path    = "../../../service-test/task"
		taskKey = "call"
		data    = `{"headers": {}}`
		server  = newServer(t)
	)

	s, validationErr, err := server.api.DeployService(serviceTar(t, path))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Id)

	require.NoError(t, server.api.StartService(s.Id))
	defer server.api.StopService(s.Id)

	_, err = server.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		ServiceID: s.Id,
		TaskKey:   taskKey,
		InputData: data,
	})
	require.NotNil(t, err)
	invalidErr := err.(*service.InvalidTaskInputError)
	require.Equal(t, taskKey, invalidErr.TaskKey)
}

func TestExecuteWithNonRunningService(t *testing.T) {
	var (
		path   = "../../../service-test/task"
		server = newServer(t)
	)

	s, validationErr, err := server.api.DeployService(serviceTar(t, path))
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Id)

	_, err = server.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		ServiceID: s.Id,
		TaskKey:   "test",
		InputData: "{}",
	})
	require.Equal(t, &api.NotRunningServiceError{ServiceID: s.Id}, err)
}

func TestExecuteWithNonExistingService(t *testing.T) {
	server := newServer(t)

	_, err := server.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		ServiceID: "service that doesnt exists",
		TaskKey:   "error",
		InputData: "{}",
	})
	require.NotNil(t, err)
}
