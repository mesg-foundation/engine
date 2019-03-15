// +build integration

package service

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func TestIntegrationStatusService(t *testing.T) {
	service, _ := FromService(&Service{
		Hash: "00",
		Name: "TestStatusService",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(newIntegrationContainer(t)))
	status, err := service.Status()
	require.NoError(t, err)
	require.Equal(t, STOPPED, status)
	dockerServices, err := service.Start()
	defer service.Stop()
	require.NoError(t, err)
	require.Equal(t, len(dockerServices), len(service.Dependencies))
	status, err = service.Status()
	require.NoError(t, err)
	require.Equal(t, RUNNING, status)
}

func TestIntegrationStatusDependency(t *testing.T) {
	service, _ := FromService(&Service{
		Hash: "00",
		Name: "TestStatusDependency",
		Dependencies: []*Dependency{
			{
				Key:   "test",
				Image: "http-server",
			},
		},
	}, ContainerOption(newIntegrationContainer(t)))
	dep := service.Dependencies[0]
	status, err := dep.Status()
	require.NoError(t, err)
	require.Equal(t, container.STOPPED, status)
	dockerServices, err := service.Start()
	require.NoError(t, err)
	require.Equal(t, len(dockerServices), len(service.Dependencies))
	status, err = dep.Status()
	require.NoError(t, err)
	require.Equal(t, container.RUNNING, status)
	service.Stop()
}
