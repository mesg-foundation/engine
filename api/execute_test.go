package api

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stvp/assert"
)

func TestNotRunningServiceError(t *testing.T) {
	e := NotRunningServiceError{ServiceID: "test"}
	assert.Equal(t, `Service "test" is not running`, e.Error())
}

func TestExecuteFunc(t *testing.T) {
	a, _, closer := newAPIAndDockerTest(t)
	defer closer()
	executor := newTaskExecutor(a)
	s, _ := service.FromService(&service.Service{
		Name: "TestExecuteFunc",
		Tasks: []*service.Task{
			{
				Key: "test",
			},
		},
	}, service.ContainerOption(a.container))
	id, err := executor.execute(s, "test", map[string]interface{}{}, []string{})
	assert.Nil(t, err)
	assert.NotNil(t, id)
}

func TestExecuteFuncInvalidTaskName(t *testing.T) {
	a, _, closer := newAPIAndDockerTest(t)
	defer closer()
	executor := newTaskExecutor(a)
	srv := &service.Service{}
	_, err := executor.execute(srv, "test", map[string]interface{}{}, []string{})
	assert.NotNil(t, err)
}

func TestCheckServiceNotRunning(t *testing.T) {
	a, _, closer := newAPIAndDockerTest(t)
	defer closer()
	executor := newTaskExecutor(a)
	err := executor.checkServiceStatus(&service.Service{Name: "TestCheckServiceNotRunning"})
	assert.NotNil(t, err)
	_, notRunningError := err.(*NotRunningServiceError)
	assert.True(t, notRunningError)
}

func TestCheckService(t *testing.T) {
	a, _, closer := newAPIAndDockerTest(t)
	defer closer()
	executor := newTaskExecutor(a)
	s, _ := service.FromService(&service.Service{
		Name: "TestCheckService",
		Dependencies: []*service.Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, service.ContainerOption(a.container))
	s.Start()
	err := executor.checkServiceStatus(s)
	assert.Nil(t, err)
}
