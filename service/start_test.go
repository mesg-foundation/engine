package service

import (
	"testing"

	"github.com/stvp/assert"
)

func TestStartService(t *testing.T) {
	service := &Service{
		Name: "TestStartService",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
		},
	}
	err := service.Start()
	assert.Nil(t, err)
	assert.Equal(t, service.IsRunning(), true)
	service.Stop()
}

func TestStartAgainService(t *testing.T) {
	service := &Service{
		Name: "TestStartAgainService",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
		},
	}
	service.Start()
	err := service.Start()
	assert.Nil(t, err)
	assert.Equal(t, service.IsRunning(), true)
	service.Stop()
}

func TestPartiallyRunningService(t *testing.T) {
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
	service.Start()
	service.Dependencies["test"].Stop(service.namespace(), "test")
	assert.Equal(t, service.IsPartiallyRunning(), true)
	err := service.Start()
	assert.Nil(t, err)
	assert.Equal(t, service.IsRunning(), true)
	service.Stop()
}
