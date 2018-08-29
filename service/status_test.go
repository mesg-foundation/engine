package service

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func TestStatusService(t *testing.T) {
	service := &Service{
		Name: "TestStatusService",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "http-server",
			},
		},
	}
	status, err := service.Status()
	require.Nil(t, err)
	require.Equal(t, STOPPED, status)
	dockerServices, err := service.Start()
	defer service.Stop()
	require.Nil(t, err)
	require.Equal(t, len(dockerServices), len(service.GetDependencies()))
	status, err = service.Status()
	require.Nil(t, err)
	require.Equal(t, RUNNING, status)
}

func TestStatusDependency(t *testing.T) {
	service := &Service{
		Name: "TestStatusDependency",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "http-server",
			},
		},
	}
	dep := service.DependenciesFromService()[0]
	status, err := dep.Status()
	require.Nil(t, err)
	require.Equal(t, container.STOPPED, status)
	dockerServices, err := service.Start()
	require.Nil(t, err)
	require.Equal(t, len(dockerServices), len(service.GetDependencies()))
	status, err = dep.Status()
	require.Nil(t, err)
	require.Equal(t, container.RUNNING, status)
	service.Stop()
}

func TestList(t *testing.T) {
	service := &Service{
		Name: "TestList",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "http-server",
			},
		},
	}
	hash := service.Hash()
	service.Start()
	defer service.Stop()
	list, err := ListRunning()
	require.Nil(t, err)
	require.Equal(t, len(list), 1)
	require.Equal(t, list[0], hash)
}

func TestListMultipleDependencies(t *testing.T) {
	service := &Service{
		Name: "TestListMultipleDependencies",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "http-server",
			},
			"test2": {
				Image: "http-server",
			},
		},
	}
	hash := service.Hash()
	service.Start()
	defer service.Stop()
	list, err := ListRunning()
	require.Nil(t, err)
	require.Equal(t, len(list), 1)
	require.Equal(t, list[0], hash)
}
