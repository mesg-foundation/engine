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
	srv := service.Service{
		Name: "TestExecute",
		Tasks: map[string]*service.Task{
			"test": &service.Task{},
		},
		Dependencies: map[string]*service.Dependency{
			"test": {
				Image: "nginx",
			},
		},
	}
	deployment, _ := serverexecute.DeployService(context.Background(), &DeployServiceRequest{
		Service: &srv,
	})
	defer services.Delete(deployment.ServiceID)
	srv.Start()
	defer srv.Stop()
	reply, err := serverexecute.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		ServiceID: deployment.ServiceID,
		TaskKey:   "test",
		InputData: "{}",
	})

	assert.Nil(t, err)
	assert.NotNil(t, reply)
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
	srv := service.Service{
		Name: "TestExecuteWithInvalidTask",
		Tasks: map[string]*service.Task{
			"test": {},
		},
		Dependencies: map[string]*service.Dependency{
			"test": {
				Image: "nginx",
			},
		},
	}
	deployment, _ := serverexecute.DeployService(context.Background(), &DeployServiceRequest{
		Service: &srv,
	})
	defer services.Delete(deployment.ServiceID)
	srv.Start()
	defer srv.Stop()
	_, err := serverexecute.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		ServiceID: deployment.ServiceID,
		TaskKey:   "error",
		InputData: "{}",
	})

	assert.NotNil(t, err)
	_, invalid := err.(*service.TaskNotFoundError)
	assert.True(t, invalid)
}

func TestExecuteWithNonRunningService(t *testing.T) {
	srv := service.Service{
		Name: "TestExecuteWithNonRunningService",
		Tasks: map[string]*service.Task{
			"test": &service.Task{},
		},
	}
	deployment, _ := serverexecute.DeployService(context.Background(), &DeployServiceRequest{
		Service: &srv,
	})
	defer services.Delete(deployment.ServiceID)
	_, err := serverexecute.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		ServiceID: deployment.ServiceID,
		TaskKey:   "test",
		InputData: "{}",
	})

	assert.NotNil(t, err)
	_, nonRunning := err.(*NotRunningServiceError)
	assert.True(t, nonRunning)
}

func TestExecuteWithNonExistingService(t *testing.T) {
	_, err := serverexecute.ExecuteTask(context.Background(), &ExecuteTaskRequest{
		ServiceID: "service that doesnt exists",
		TaskKey:   "error",
		InputData: "{}",
	})

	assert.NotNil(t, err)
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
