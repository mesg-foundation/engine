// +build integration

package dockermanager

import (
	"testing"

	"github.com/mesg-foundation/core/service"
	"github.com/stretchr/testify/require"
)

func TestIntegrationLogs(t *testing.T) {
	var (
		service = &service.Service{
			Hash: "1",
			Name: "TestLogs",
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
	readers, err := m.Logs(service)
	require.NoError(t, err)
	require.Equal(t, 2, len(readers))
}

func TestIntegrationLogsOnlyOneDependency(t *testing.T) {
	var (
		service = &service.Service{
			Hash: "1",
			Name: "TestLogsOnlyOneDependency",
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
	readers, err := m.Logs(service, "test2")
	require.NoError(t, err)
	require.Equal(t, 1, len(readers))
}
