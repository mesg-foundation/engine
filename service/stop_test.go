package service

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func TestStopRunningService(t *testing.T) {
	service := &Service{
		Name: "TestStopRunningService",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "http-server",
			},
		},
	}
	service.Start()
	err := service.Stop()
	require.Nil(t, err)
	status, _ := service.Status()
	require.Equal(t, STOPPED, status)
}

func TestStopNonRunningService(t *testing.T) {
	service := &Service{
		Name: "TestStopNonRunningService",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "http-server",
			},
		},
	}
	err := service.Stop()
	require.Nil(t, err)
	status, _ := service.Status()
	require.Equal(t, STOPPED, status)
}

func TestStopDependency(t *testing.T) {
	service := &Service{
		Name: "TestStartDependency",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "http-server",
			},
		},
	}
	networkID, err := defaultContainer.CreateNetwork(service.namespace())
	defer defaultContainer.DeleteNetwork(service.namespace())
	dep := service.DependenciesFromService()[0]
	dep.Start(networkID)
	err = dep.Stop()
	require.Nil(t, err)
	status, _ := dep.Status()
	require.Equal(t, container.STOPPED, status)
}

func TestNetworkDeleted(t *testing.T) {
	service := &Service{
		Name: "TestNetworkDeleted",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "http-server",
			},
		},
	}
	service.Start()
	service.Stop()
	n, err := defaultContainer.FindNetwork(service.namespace())
	require.Empty(t, n)
	require.NotNil(t, err)
}
