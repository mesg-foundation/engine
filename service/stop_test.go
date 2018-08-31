package service

import (
	"testing"
	"time"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func TestStopRunningService(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestStopRunningService",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "nginx:stable-alpine",
			},
		},
	}, ContainerOption(defaultContainer))
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
				Image: "nginx:stable-alpine",
			},
		},
	}, ContainerOption(defaultContainer))
	err := service.Stop()
	require.Nil(t, err)
	status, _ := service.Status()
	require.Equal(t, STOPPED, status)
}

func TestStopDependency(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestStartDependency",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "nginx:stable-alpine",
			},
		},
	}, ContainerOption(defaultContainer))
	networkID, err := defaultContainer.CreateNetwork(service.namespace())
	defer defaultContainer.DeleteNetwork(service.namespace())
	dep := service.Dependencies[0]
	dep.Start(networkID)
	err = dep.Stop()
	require.Nil(t, err)
	status, _ := dep.Status()
	require.Equal(t, container.STOPPED, status)
}

func TestNetworkDeleted(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestNetworkDeleted",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "nginx:stable-alpine",
			},
		},
	}, ContainerOption(defaultContainer))
	service.Start()
	service.Stop()
	time.Sleep(5 * time.Second)
	_, err := defaultContainer.FindNetwork(service.namespace())
	require.NotNil(t, err)
}
