package service

import (
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/stvp/assert"
)

func TestStatusService(t *testing.T) {
	service := &Service{
		Name: "TestStatusService",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
		},
	}
	status, err := service.Status()
	assert.Nil(t, err)
	assert.Equal(t, STOPPED, status)
	dockerServices, err := service.Start()
	defer service.Stop()
	assert.Nil(t, err)
	assert.Equal(t, len(dockerServices), len(service.GetDependencies()))
	status, err = service.Status()
	assert.Nil(t, err)
	assert.Equal(t, RUNNING, status)
}

func TestStatusDependency(t *testing.T) {
	service := &Service{
		Name: "TestStatusDependency",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
		},
	}
	dep := service.DependenciesFromService()[0]
	status, err := dep.Status()
	assert.Nil(t, err)
	assert.Equal(t, container.STOPPED, status)
	dockerServices, err := service.Start()
	assert.Nil(t, err)
	assert.Equal(t, len(dockerServices), len(service.GetDependencies()))
	status, err = dep.Status()
	assert.Nil(t, err)
	assert.Equal(t, container.RUNNING, status)
	service.Stop()
}

func TestList(t *testing.T) {
	service := &Service{
		Name: "TestList",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
		},
	}
	hash := service.Hash()
	service.Start()
	defer service.Stop()
	list, err := ListRunning()
	assert.Nil(t, err)
	assert.Equal(t, len(list), 1)
	assert.Equal(t, list[0], hash)
}

func TestListMultipleDependencies(t *testing.T) {
	service := &Service{
		Name: "TestListMultipleDependencies",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
			"test2": &Dependency{
				Image: "nginx",
			},
		},
	}
	hash := service.Hash()
	service.Start()
	defer service.Stop()
	list, err := ListRunning()
	assert.Nil(t, err)
	assert.Equal(t, len(list), 1)
	assert.Equal(t, list[0], hash)
}
