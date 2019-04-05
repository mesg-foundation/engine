// +build integration

package service

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func TestIntegrationStopRunningService(t *testing.T) {
	service := &Service{
		Hash: "1",
		Name: "TestStopRunningService",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}
	c := newIntegrationContainer(t)
	service.Start(c)
	err := service.Stop(c)
	require.NoError(t, err)
	status, _ := service.Status(c)
	require.Equal(t, STOPPED, status)
}

func TestIntegrationStopNonRunningService(t *testing.T) {
	service := &Service{
		Hash: "1",
		Name: "TestStopNonRunningService",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}
	c := newIntegrationContainer(t)
	err := service.Stop(c)
	require.NoError(t, err)
	status, _ := service.Status(c)
	require.Equal(t, STOPPED, status)
}

func TestIntegrationStopDependency(t *testing.T) {
	service := &Service{
		Hash: "1",
		Name: "TestStopDependency",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}
	c := newIntegrationContainer(t)
	networkID, err := c.CreateNetwork(service.namespace())
	require.NoError(t, err)
	defer c.DeleteNetwork(service.namespace())
	dep := service.Dependencies[0]
	dep.Start(c, service, networkID)
	err = dep.Stop(c, service)
	require.NoError(t, err)
	status, _ := dep.Status(c, service)
	require.Equal(t, container.STOPPED, status)
}

func TestIntegrationNetworkDeleted(t *testing.T) {
	service := &Service{
		Hash: "1",
		Name: "TestNetworkDeleted",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}
	c := newIntegrationContainer(t)
	service.Start(c)
	service.Stop(c)
	n, err := c.FindNetwork(service.namespace())
	require.Empty(t, n)
	require.Error(t, err)
}
