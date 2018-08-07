package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

var serverexecute = new(Server)

func TestExecute(t *testing.T) {
	deployment, _ := serverexecute.DeployService(context.Background(), &DeployServiceRequest{
		Service: &service.Service{
			Name: "TestExecute",
			Tasks: map[string]*service.Task{
				"test": {},
			},
		},
	})
	reply, err := serverexecute.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		ServiceID: deployment.ServiceID,
		TaskKey:   "test",
		InputData: "{}",
	})

	assert.Nil(t, err)
	assert.NotNil(t, reply)
	services.Delete(deployment.ServiceID)
}

func TestExecuteWithInvalidJSON(t *testing.T) {
	deployment, _ := serverexecute.DeployService(context.Background(), &DeployServiceRequest{
		Service: &service.Service{
			Name: "TestExecuteWithInvalidJSON",
			Tasks: map[string]*service.Task{
				"test": {},
			},
		},
	})
	_, err := serverexecute.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		ServiceID: deployment.ServiceID,
		TaskKey:   "test",
		InputData: "",
	})

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "unexpected end of JSON input")
	services.Delete(deployment.ServiceID)
}

func TestExecuteWithInvalidTask(t *testing.T) {
	deployment, _ := serverexecute.DeployService(context.Background(), &DeployServiceRequest{
		Service: &service.Service{
			Name: "TestExecuteWithInvalidJSON",
			Tasks: map[string]*service.Task{
				"test": {},
			},
		},
	})
	_, err := serverexecute.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		ServiceID: deployment.ServiceID,
		TaskKey:   "error",
		InputData: "{}",
	})

	assert.NotNil(t, err)
	_, invalid := err.(*service.TaskNotFoundError)
	assert.True(t, invalid)
	services.Delete(deployment.ServiceID)
}

func TestExecuteWithNonExistingService(t *testing.T) {
	_, err := serverexecute.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		ServiceID: "service that doesnt exists",
		TaskKey:   "error",
		InputData: "{}",
	})

	assert.NotNil(t, err)
}
