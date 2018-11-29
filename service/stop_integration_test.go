// +build integration

package service

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func TestIntegrationStopRunningService(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestStopRunningService",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(newIntegrationContainer(t)))
	service.Start()
	err := service.Stop()
	require.NoError(t, err)
	status, _ := service.Status()
	require.Equal(t, STOPPED, status)
}

func TestIntegrationStopNonRunningService(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestStopNonRunningService",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(newIntegrationContainer(t)))
	err := service.Stop()
	require.NoError(t, err)
	status, _ := service.Status()
	require.Equal(t, STOPPED, status)
}

func TestIntegrationStopDependency(t *testing.T) {
	c := newIntegrationContainer(t)
	service, _ := FromService(&Service{
		Name: "TestStopDependency",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(c))
	networkID, err := c.CreateNetwork(service.namespace())
	require.NoError(t, err)
	defer c.DeleteNetwork(service.namespace(), container.EventDestroy)
	dep := service.Dependencies[0]
	dep.Start(networkID)
	err = dep.Stop()
	require.NoError(t, err)
	status, _ := dep.Status()
	require.Equal(t, container.STOPPED, status)
}

func TestIntegrationNetworkDeleted(t *testing.T) {
	c := newIntegrationContainer(t)
	service, _ := FromService(&Service{
		Name: "TestNetworkDeleted",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(c))
	service.Start()
	service.Stop()
	n, err := c.FindNetwork(service.namespace())
	require.Empty(t, n)
	require.Error(t, err)
}
