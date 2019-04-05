// +build integration

package service

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/stretchr/testify/require"
)

func TestIntegrationStatusService(t *testing.T) {
	var (
		service = &Service{
			Hash: "00",
			Name: "TestStatusService",
			Dependencies: []*Dependency{
				{
					Key:   "test",
					Image: "http-server",
				},
			},
		}
		c = newIntegrationContainer(t)
	)
	status, err := service.Status(c)
	require.NoError(t, err)
	require.Equal(t, STOPPED, status)
	dockerServices, err := service.Start(c)
	defer service.Stop(c)
	require.NoError(t, err)
	require.Equal(t, len(dockerServices), len(service.Dependencies))
	status, err = service.Status(c)
	require.NoError(t, err)
	require.Equal(t, RUNNING, status)
}

func TestIntegrationStatusDependency(t *testing.T) {
	var (
		service = &Service{
			Hash: "00",
			Name: "TestStatusDependency",
			Dependencies: []*Dependency{
				{
					Key:   "test",
					Image: "http-server",
				},
			},
		}
		c = newIntegrationContainer(t)
	)
	dep := service.Dependencies[0]
	status, err := dep.Status(c, service)
	require.NoError(t, err)
	require.Equal(t, container.STOPPED, status)
	dockerServices, err := service.Start(c)
	require.NoError(t, err)
	require.Equal(t, len(dockerServices), len(service.Dependencies))
	status, err = dep.Status(c, service)
	require.NoError(t, err)
	require.Equal(t, container.RUNNING, status)
	service.Stop(c)
}

func TestIntegrationListRunning(t *testing.T) {
	var (
		service = &Service{
			Hash: "00",
			Name: "TestList",
			Dependencies: []*Dependency{
				{
					Key:   "test",
					Image: "http-server",
				},
			},
		}
		c = newIntegrationContainer(t)
	)
	service.Start(c)
	defer service.Stop(c)
	list, err := ListRunning()
	require.NoError(t, err)
	require.Equal(t, len(list), 1)
	require.Equal(t, list[0], service.Hash)
}

func TestIntegrationListRunningMultipleDependencies(t *testing.T) {
	var (
		service = &Service{
			Hash: "00",
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
		}
		c = newIntegrationContainer(t)
	)
	service.Start(c)
	defer service.Stop(c)
	list, err := ListRunning()
	require.NoError(t, err)
	require.Equal(t, len(list), 1)
	require.Equal(t, list[0], service.Hash)
}
