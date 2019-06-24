// +build integration

package dockermanager

import (
	"testing"

	"github.com/mesg-foundation/core/hash"
	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestIntegrationStatusService(t *testing.T) {
	var (
		s = &service.Service{
			Hash: hash.Int(1),
			Name: "TestStatusService",
			Dependencies: []*service.Dependency{
				{
					Key:   "test",
					Image: "http-server",
				},
			},
		}
		c = newIntegrationContainer(t)
		m = New(c)
	)

	status, err := m.Status(s)
	require.NoError(t, err)
	require.Equal(t, service.STOPPED, status)
	dockerServices, err := m.Start(s)
	defer m.Stop(s)
	require.NoError(t, err)
	require.Equal(t, len(dockerServices), len(s.Dependencies))
	status, err = m.Status(s)
	require.NoError(t, err)
	require.Equal(t, service.RUNNING, status)
}

func TestIntegrationListRunning(t *testing.T) {
	var (
		service = &service.Service{
			Hash: hash.Int(1),
			Name: "TestList",
			Dependencies: []*service.Dependency{
				{
					Key:   "test",
					Image: "http-server",
				},
			},
		}
		c = newIntegrationContainer(t)
		m = New(c)
	)

	m.Start(service)
	defer m.Stop(service)
	list, err := ListRunning()
	require.NoError(t, err)
	require.Equal(t, len(list), 1)
	require.Equal(t, list[0], service.Hash.String())
}

func TestIntegrationListRunningMultipleDependencies(t *testing.T) {
	var (
		service = &service.Service{
			Hash: hash.Int(1),
			Name: "TestListMultipleDependencies",
			Dependencies: []*service.Dependency{
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
		m = New(c)
	)

	m.Start(service)
	defer m.Stop(service)
	list, err := ListRunning()
	require.NoError(t, err)
	require.Equal(t, len(list), 1)
	require.Equal(t, list[0], service.Hash.String())
}
