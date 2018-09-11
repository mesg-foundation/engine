package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogs(t *testing.T) {
	service := &Service{
		Name: "TestLogs",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "nginx:stable-alpine",
			},
			"test2": {
				Image: "nginx:stable-alpine",
			},
		},
	}
	service.Start()
	defer service.Stop()
	readers, err := service.Logs("*")
	require.Nil(t, err)
	require.Equal(t, 2, len(readers))
}

func TestLogsOnlyOneDependency(t *testing.T) {
	service := &Service{
		Name: "TestLogsOnlyOneDependency",
		Dependencies: map[string]*Dependency{
			"test": {
				Image: "nginx:stable-alpine",
			},
			"test2": {
				Image: "nginx:stable-alpine",
			},
		},
	}
	service.Start()
	defer service.Stop()
	readers, err := service.Logs("test2")
	require.Nil(t, err)
	require.Equal(t, 1, len(readers))
}
