package service

import (
	"testing"

	"github.com/stvp/assert"
)

func TestDependenciesFromService(t *testing.T) {
	service := &Service{
		Name: "TestPartiallyRunningService",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
			"test2": &Dependency{
				Image: "nginx",
			},
		},
	}
	deps := service.DependenciesFromService()
	assert.Equal(t, 2, len(deps))
	assert.Equal(t, "test", deps[0].Name)
	assert.Equal(t, "TestPartiallyRunningService", deps[0].Service.Name)
	assert.Equal(t, "test2", deps[1].Name)
}
