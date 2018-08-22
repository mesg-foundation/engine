package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

var serverexecute = new(Server)

func TestExecute(t *testing.T) {
	var (
		url     = "https://github.com/mesg-foundation/service-webhook"
		taskKey = "call"
		data    = `{"url": "https://mesg.tech", "data": {}, "headers": {}}`
	)

	server := newServer(t)
	stream := newTestDeployStream(url)

	server.DeployService(stream)
	defer services.Delete(stream.serviceID)

	serverexecute.StartService(context.Background(), &StartServiceRequest{
		ServiceID: stream.serviceID,
	})
	defer serverexecute.StopService(context.Background(), &StopServiceRequest{
		ServiceID: stream.serviceID,
	})

	reply, err := serverexecute.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		ServiceID: stream.serviceID,
		TaskKey:   taskKey,
		InputData: data,
	})

	require.Nil(t, err)
	require.NotEqual(t, "", reply.ExecutionID)
}

func TestExecuteWithInvalidJSON(t *testing.T) {
	url := "https://github.com/mesg-foundation/service-webhook"

	server := newServer(t)
	stream := newTestDeployStream(url)
	server.DeployService(stream)

	_, err := serverexecute.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		ServiceID: stream.serviceID,
		TaskKey:   "test",
		InputData: "",
	})
	require.NotNil(t, err)
	require.Equal(t, err.Error(), "unexpected end of JSON input")
	services.Delete(stream.serviceID)
}

func TestExecuteWithInvalidTask(t *testing.T) {
	url := "https://github.com/mesg-foundation/service-webhook"

	server := newServer(t)
	stream := newTestDeployStream(url)

	server.DeployService(stream)
	defer services.Delete(stream.serviceID)

	serverexecute.StartService(context.Background(), &StartServiceRequest{
		ServiceID: stream.serviceID,
	})
	defer serverexecute.StopService(context.Background(), &StopServiceRequest{
		ServiceID: stream.serviceID,
	})

	_, err := serverexecute.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		ServiceID: stream.serviceID,
		TaskKey:   "error",
		InputData: "{}",
	})

	require.Error(t, err)
	require.IsType(t, (*service.TaskNotFoundError)(nil), err)
}

func TestExecuteWithNonRunningService(t *testing.T) {
	url := "https://github.com/mesg-foundation/service-webhook"

	server := newServer(t)
	stream := newTestDeployStream(url)

	server.DeployService(stream)
	defer services.Delete(stream.serviceID)

	_, err := serverexecute.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		ServiceID: stream.serviceID,
		TaskKey:   "test",
		InputData: "{}",
	})

	require.Equal(t, &api.NotRunningServiceError{ServiceID: stream.serviceID}, err)
}

func TestExecuteWithNonExistingService(t *testing.T) {
	_, err := serverexecute.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		ServiceID: "service that doesnt exists",
		TaskKey:   "error",
		InputData: "{}",
	})

	require.NotNil(t, err)
}
