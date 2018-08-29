package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDependenciesFromService(t *testing.T) {
	service := &Service{
		Name: "TestPartiallyRunningService",
		Dependencies: map[string]*Dependency{
			"testa": {
				Image: "nginx:stable-alpine",
			},
			"testb": {
				Image: "nginx:stable-alpine",
			},
		},
	}
	deps := service.DependenciesFromService()
	require.Equal(t, 2, len(deps))
	require.Equal(t, "testa", deps[0].Name)
	require.Equal(t, "TestPartiallyRunningService", deps[0].Service.Name)
	require.Equal(t, "testb", deps[1].Name)
}
