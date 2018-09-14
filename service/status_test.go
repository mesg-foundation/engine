package service

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func TestStatusService(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestStatusService",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(defaultContainer))
	status, err := service.Status()
	require.Nil(t, err)
	require.Equal(t, STOPPED, status)
	dockerServices, err := service.Start()
	defer service.Stop()
	require.Nil(t, err)
	require.Equal(t, len(dockerServices), len(service.Dependencies))
	status, err = service.Status()
	require.Nil(t, err)
	require.Equal(t, RUNNING, status)
}

func TestStatusDependency(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestStatusDependency",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(defaultContainer))
	dep := service.Dependencies[0]
	status, err := dep.Status()
	require.Nil(t, err)
	require.Equal(t, container.STOPPED, status)
	dockerServices, err := service.Start()
	require.Nil(t, err)
	require.Equal(t, len(dockerServices), len(service.Dependencies))
	status, err = dep.Status()
	require.Nil(t, err)
	require.Equal(t, container.RUNNING, status)
	service.Stop()
}

func TestList(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestList",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(defaultContainer))
	service.Start()
	defer service.Stop()
	list, err := ListRunning()
	require.Nil(t, err)
	require.Equal(t, len(list), 1)
	require.Equal(t, list[0], service.ID)
}

func TestListMultipleDependencies(t *testing.T) {
	service, _ := FromService(&Service{
		Name: "TestListMultipleDependencies",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
			{
				Key:   "test2",
				Image: "http-server",
			},
		},
	}, ContainerOption(defaultContainer))
	service.Start()
	defer service.Stop()
	list, err := ListRunning()
	require.Nil(t, err)
	require.Equal(t, len(list), 1)
	require.Equal(t, list[0], service.ID)
}
