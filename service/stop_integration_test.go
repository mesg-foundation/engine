// +build integration

package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationStopRunningService(t *testing.T) {
	var (
		service = &Service{
			Hash: "1",
			Name: "TestStopRunningService",
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
	err := service.Stop(c)
	require.NoError(t, err)
	status, _ := service.Status(c)
	require.Equal(t, STOPPED, status)
}

func TestIntegrationStopNonRunningService(t *testing.T) {
	var (
		service = &Service{
			Hash: "1",
			Name: "TestStopNonRunningService",
			Dependencies: []*Dependency{
				{
					Key:   "test",
					Image: "http-server",
				},
			},
		}
		c = newIntegrationContainer(t)
	)

	err := service.Stop(c)
	require.NoError(t, err)
	status, _ := service.Status(c)
	require.Equal(t, STOPPED, status)
}

func TestIntegrationNetworkDeleted(t *testing.T) {
	var (
		service = &Service{
			Hash: "1",
			Name: "TestNetworkDeleted",
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
	service.Stop(c)
	n, err := c.FindNetwork(service.namespace())
	require.Empty(t, n)
	require.Error(t, err)
}
