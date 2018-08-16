package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
	"github.com/stvp/assert"
)

var serverexecute = new(Server)

func TestExecute(t *testing.T) {
	url := "https://github.com/mesg-foundation/service-webhook"
	taskKey := "call"
	data := `{"url": "https://mesg.tech", "data": {}, "headers": {}}`

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

	require.NotNil(t, err)
	_, nonRunning := err.(*NotRunningServiceError)
	require.True(t, nonRunning)
}

func TestExecuteWithNonExistingService(t *testing.T) {
	_, err := serverexecute.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		ServiceID: "service that doesnt exists",
		TaskKey:   "error",
		InputData: "{}",
	})

	require.NotNil(t, err)
}

func TestExecuteFunc(t *testing.T) {
	srv := &service.Service{
		Name: "TestExecuteFunc",
		Tasks: map[string]*service.Task{
			"test": {},
		},
	}
	id, err := execute(srv, "test", map[string]interface{}{})
	assert.Nil(t, err)
	assert.NotNil(t, id)
}

func TestExecuteFuncInvalidTaskName(t *testing.T) {
	srv := &service.Service{}
	_, err := execute(srv, "test", map[string]interface{}{})
	assert.NotNil(t, err)
}

func TestGetData(t *testing.T) {
	inputs, err := getData(&ExecuteTaskRequest{
		InputData: "{\"foo\":\"bar\"}",
	})
	assert.Nil(t, err)
	assert.Equal(t, "bar", inputs["foo"])
}

func TestGetDataInvalid(t *testing.T) {
	_, err := getData(&ExecuteTaskRequest{
		InputData: "",
	})
	assert.NotNil(t, err)
}

func TestCheckServiceNotRunning(t *testing.T) {
	err := checkServiceStatus(&service.Service{Name: "TestCheckServiceNotRunning"})
	assert.NotNil(t, err)
	_, notRunningError := err.(*NotRunningServiceError)
	assert.True(t, notRunningError)
}

func TestCheckService(t *testing.T) {
	srv := service.Service{
		Name: "TestCheckService",
		Dependencies: map[string]*service.Dependency{
			"test": {
				Image: "nginx",
			},
		},
	}
	srv.Start()
	defer srv.Stop()
	err := checkServiceStatus(&srv)
	assert.Nil(t, err)
}
