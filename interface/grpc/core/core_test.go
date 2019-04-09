package core

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestGetService(t *testing.T) {
	server, closer := newServer(t)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, taskServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Hash, false)

	reply, err := server.GetService(context.Background(), &coreapi.GetServiceRequest{
		ServiceID: s.Hash,
	})
	require.NoError(t, err)
	require.NotNil(t, reply)
	require.Equal(t, reply.Service.Definition.Name, "Task")
}

func TestListServices(t *testing.T) {
	url := "https://github.com/mesg-foundation/service-webhook"
	server, closer := newServer(t)
	defer closer()

	stream := newTestDeployStream(url)
	require.NoError(t, server.DeployService(stream))
	defer server.api.DeleteService(stream.hash, false)

	reply, err := server.ListServices(context.Background(), &coreapi.ListServicesRequest{})
	require.NoError(t, err)

	services, err := server.api.ListServices()
	require.NoError(t, err)

	apiProtoServices := toProtoServices(services)

	require.Len(t, apiProtoServices, 1)
	require.Equal(t, reply.Services[0].Definition.Hash, apiProtoServices[0].Hash)
}

func TestStartService(t *testing.T) {
	server, closer := newServer(t)
	defer closer()

	// we use a test service without tasks definition here otherwise we need to
	// spin up the gRPC server in order to prevent service exit with a failure
	// because it'll try to listen for tasks.
	s, validationErr, err := server.api.DeployService(serviceTar(t, eventServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Hash, false)

	_, err = server.StartService(context.Background(), &coreapi.StartServiceRequest{
		ServiceID: s.Hash,
	})
	require.NoError(t, err)
	defer server.api.StopService(s.Hash)

	status, err := server.api.Status(s)
	require.NoError(t, err)
	require.Equal(t, service.RUNNING, status)
}

func TestStopService(t *testing.T) {
	server, closer := newServer(t)
	defer closer()

	// we use a test service without tasks definition here otherwise we need to
	// spin up the gRPC server in order to prevent service exit with a failure
	// because it'll try to listen for tasks.
	s, validationErr, err := server.api.DeployService(serviceTar(t, eventServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Hash, false)

	require.NoError(t, server.api.StartService(s.Hash))

	reply, err := server.StopService(context.Background(), &coreapi.StopServiceRequest{
		ServiceID: s.Hash,
	})
	require.NoError(t, err)

	status, err := server.api.Status(s)
	require.NoError(t, err)
	require.Equal(t, service.STOPPED, status)
	require.NoError(t, err)
	require.NotNil(t, reply)
}

func TestDeleteService(t *testing.T) {
	var (
		path           = filepath.Join("..", "..", "..", "service-test", "task")
		server, closer = newServer(t)
	)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, path), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)

	reply, err := server.DeleteService(context.Background(), &coreapi.DeleteServiceRequest{
		ServiceID: s.Hash,
	})
	require.NoError(t, err)
	require.NotNil(t, reply)

	_, err = server.api.GetService(s.Hash)
	require.Error(t, err)
}

func TestExecute(t *testing.T) {
	var (
		taskKey        = "call"
		data           = `{"url": "https://mesg.com", "data": {}, "headers": {}}`
		server, closer = newServer(t)
	)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, taskServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Hash, false)

	require.NoError(t, server.api.StartService(s.Hash))
	defer server.api.StopService(s.Hash)

	reply, err := server.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID: s.Sid,
		TaskKey:   taskKey,
		InputData: data,
	})
	require.NoError(t, err)
	require.NotEqual(t, "", reply.ExecutionID)
}

func TestExecuteWithInvalidJSON(t *testing.T) {
	server, closer := newServer(t)
	defer closer()

	s, validationErr, err := server.api.DeployService(serviceTar(t, taskServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Hash, false)

	_, err = server.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID: s.Sid,
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

	s, validationErr, err := server.api.DeployService(serviceTar(t, taskServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Hash, false)

	require.NoError(t, server.api.StartService(s.Hash))
	defer server.api.StopService(s.Hash)

	_, err = server.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID: s.Sid,
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

	s, validationErr, err := server.api.DeployService(serviceTar(t, taskServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Hash, false)

	require.NoError(t, server.api.StartService(s.Hash))
	defer server.api.StopService(s.Hash)

	_, err = server.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID: s.Sid,
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

	s, validationErr, err := server.api.DeployService(serviceTar(t, taskServicePath), nil)
	require.Zero(t, validationErr)
	require.NoError(t, err)
	defer server.api.DeleteService(s.Hash, false)

	_, err = server.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID: s.Sid,
		TaskKey:   "test",
		InputData: "{}",
	})
	require.Equal(t, &api.NotRunningServiceError{ServiceID: s.Sid}, err)
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
	require.Equal(t, reply.Address, c.Client.Address)
	require.Equal(t, reply.Image, c.Core.Image)
}
