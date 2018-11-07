package api

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestNotRunningServiceError(t *testing.T) {
	e := NotRunningServiceError{ServiceID: "test"}
	require.Equal(t, `Service "test" is not running`, e.Error())
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
	id, err := executor.execute(s, "xxx", "test", map[string]interface{}{}, []string{})
	require.NoError(t, err)
	require.NotNil(t, id)
}

func TestExecuteFuncInvalidTaskName(t *testing.T) {
	a, _, closer := newAPIAndDockerTest(t)
	defer closer()
	executor := newTaskExecutor(a)
	srv := &service.Service{}
	_, err := executor.execute(srv, "xxx", "test", map[string]interface{}{}, []string{})
	require.Error(t, err)
}

func TestCheckServiceNotRunning(t *testing.T) {
	a, _, closer := newAPIAndDockerTest(t)
	defer closer()
	executor := newTaskExecutor(a)
	err := executor.checkServiceStatus(&service.Service{Name: "TestCheckServiceNotRunning"})
	require.Error(t, err)
	_, notRunningError := err.(*NotRunningServiceError)
	require.True(t, notRunningError)
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
	require.NoError(t, err)
}
