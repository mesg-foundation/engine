// +build integration

package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationLogs(t *testing.T) {
	var (
		service = &Service{
			Hash: "1",
			Name: "TestLogs",
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
	readers, err := service.Logs(c)
	require.NoError(t, err)
	require.Equal(t, 2, len(readers))
}

func TestIntegrationLogsOnlyOneDependency(t *testing.T) {
	var (
		service = &Service{
			Hash: "1",
			Name: "TestLogsOnlyOneDependency",
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
	readers, err := service.Logs(c, "test2")
	require.NoError(t, err)
	require.Equal(t, 1, len(readers))
}
