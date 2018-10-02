// +build integration

package service

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func TestStopRunningService(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestStopRunningService",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(newContainer(t)))
	service.Start()
	err := service.Stop()
	require.Nil(t, err)
	status, _ := service.Status()
	require.Equal(t, STOPPED, status)
}

func TestStopNonRunningService(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestStopNonRunningService",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(newContainer(t)))
	err := service.Stop()
	require.Nil(t, err)
	status, _ := service.Status()
	require.Equal(t, STOPPED, status)
}

func TestStopDependency(t *testing.T) {
	c := newContainer(t)
	service, _ := FromService(&Service{
		Name: "TestStartDependency",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(c))
	networkID, err := c.CreateNetwork(service.namespace())
	defer c.DeleteNetwork(service.namespace(), container.EventDestroy)
	dep := service.Dependencies[0]
	dep.Start(networkID)
	err = dep.Stop()
	require.Nil(t, err)
	status, _ := dep.Status()
	require.Equal(t, container.STOPPED, status)
}

func TestNetworkDeleted(t *testing.T) {
	c := newContainer(t)
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
	require.NotNil(t, err)
}
