// +build integration

package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	executionsdk "github.com/mesg-foundation/core/sdk/execution"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	var (
		taskKey        = "call"
		data           = `{"url": "https://mesg.com", "data": {}, "headers": {}}`
		server, closer = newServer(t)
	)
	defer closer()

	s, validationErr, err := server.sdk.DeployService(serviceTar(t, taskServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.sdk.DeleteService(s.Hash, false)

	require.NoError(t, server.sdk.StartService(s.Hash))
	defer server.sdk.StopService(s.Hash)

	reply, err := server.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID: s.Hash.String(),
		TaskKey:   taskKey,
		InputData: data,
	})
	require.NoError(t, err)
	require.NotEqual(t, "", reply.ExecutionHash)
}

func TestExecuteWithInvalidJSON(t *testing.T) {
	server, closer := newServer(t)
	defer closer()

	s, validationErr, err := server.sdk.DeployService(serviceTar(t, taskServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.sdk.DeleteService(s.Hash, false)

	_, err = server.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID: s.Hash.String(),
		TaskKey:   "test",
		InputData: "",
	})
	require.Error(t, err)
	require.Equal(t, err.Error(), "cannot parse execution's inputs (JSON format): unexpected end of JSON input")
}

func TestExecuteWithInvalidTask(t *testing.T) {
	var (
		taskKey        = "error"
		server, closer = newServer(t)
	)
	defer closer()

	s, validationErr, err := server.sdk.DeployService(serviceTar(t, taskServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.sdk.DeleteService(s.Hash, false)

	require.NoError(t, server.sdk.StartService(s.Hash))
	defer server.sdk.StopService(s.Hash)

	_, err = server.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID: s.Hash.String(),
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

	s, validationErr, err := server.sdk.DeployService(serviceTar(t, taskServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.sdk.DeleteService(s.Hash, false)

	require.NoError(t, server.sdk.StartService(s.Hash))
	defer server.sdk.StopService(s.Hash)

	_, err = server.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID: s.Hash.String(),
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

	s, validationErr, err := server.sdk.DeployService(serviceTar(t, taskServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.sdk.DeleteService(s.Hash, false)

	_, err = server.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID: s.Hash.String(),
		TaskKey:   "test",
		InputData: "{}",
	})
	require.Equal(t, &executionsdk.NotRunningServiceError{ServiceID: s.Hash.String()}, err)
}

func TestExecuteWithNonExistingService(t *testing.T) {
	server, closer := newServer(t)
	defer closer()

	_, err := server.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID: "-",
		TaskKey:   "error",
		InputData: "{}",
	})
	require.Error(t, err)
}

func TestInfo(t *testing.T) {
	server, closer := newServer(t)
	defer closer()

	c, err := config.Global()
	require.NoError(t, err)
	reply, err := server.Info(context.Background(), &coreapi.InfoRequest{})
	require.NoError(t, err)
	require.NotNil(t, reply)
	for i, s := range reply.Services {
		require.Equal(t, s.Sid, c.Services()[i].Sid)
	}
}
