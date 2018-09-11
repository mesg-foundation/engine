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
	a, _ := newAPIAndDockerTest(t)
	executor := newTaskExecutor(a)
	s := &service.Service{
		Name: "TestExecuteFunc",
		Tasks: map[string]*service.Task{
			"test": {},
		},
	}
	id, err := executor.execute(s, "test", map[string]interface{}{}, []string{})
	assert.Nil(t, err)
	assert.NotNil(t, id)
}

func TestExecuteFuncInvalidTaskName(t *testing.T) {
	a, _ := newAPIAndDockerTest(t)
	executor := newTaskExecutor(a)
	srv := &service.Service{}
	_, err := executor.execute(srv, "test", map[string]interface{}{}, []string{})
	assert.NotNil(t, err)
}

func TestCheckServiceNotRunning(t *testing.T) {
	a, _ := newAPIAndDockerTest(t)
	executor := newTaskExecutor(a)
	err := executor.checkServiceStatus(&service.Service{Name: "TestCheckServiceNotRunning"})
	assert.NotNil(t, err)
	_, notRunningError := err.(*NotRunningServiceError)
	assert.True(t, notRunningError)
}

func TestCheckService(t *testing.T) {
	a, _ := newAPIAndDockerTest(t)
	executor := newTaskExecutor(a)
	s := service.Service{
		Name: "TestCheckService",
		Dependencies: map[string]*service.Dependency{
			"test": {
				Image: "nginx:stable-alpine",
			},
		},
	}
	s.Start()
	defer s.Stop()
	err := executor.checkServiceStatus(&s)
	assert.Nil(t, err)
}
