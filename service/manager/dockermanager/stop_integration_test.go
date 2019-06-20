// +build integration

package dockermanager

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestIntegrationStopRunningService(t *testing.T) {
	var (
		s = &service.Service{
			Hash: []byte{0},
			Name: "TestStopRunningService",
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

	m.Start(s)
	err := m.Stop(s)
	require.NoError(t, err)
	status, _ := m.Status(s)
	require.Equal(t, service.STOPPED, status)
}

func TestIntegrationStopNonRunningService(t *testing.T) {
	var (
		s = &service.Service{
			Hash: []byte{0},
			Name: "TestStopNonRunningService",
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

	err := m.Stop(s)
	require.NoError(t, err)
	status, _ := m.Status(s)
	require.Equal(t, service.STOPPED, status)
}

func TestIntegrationNetworkDeleted(t *testing.T) {
	var (
		service = &service.Service{
			Hash: []byte{0},
			Name: "TestNetworkDeleted",
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
	m.Stop(service)
	n, err := c.FindNetwork(serviceNamespace(service.Hash))
	require.Empty(t, n)
	require.Error(t, err)
}
