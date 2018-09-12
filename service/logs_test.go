package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogs(t *testing.T) {
	service, _ := FromService(&Service{
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
	}, ContainerOption(defaultContainer))
	service.Start()
	defer service.Stop()
	readers, err := service.Logs()
	require.Nil(t, err)
	require.Equal(t, 2, len(readers))
}

func TestLogsOnlyOneDependency(t *testing.T) {
	service, _ := FromService(&Service{
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
	}, ContainerOption(defaultContainer))
	service.Start()
	defer service.Stop()
	readers, err := service.Logs("test2")
	require.Nil(t, err)
	require.Equal(t, 1, len(readers))
}
