package service

import (
	"testing"

	"github.com/stvp/assert"
)

func TestStatusRunning(t *testing.T) {
	service := &Service{
		Name: "TestStatusRunning",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
		},
	}
	dockerServices, err := service.Start()
	defer service.Stop()
	assert.Nil(t, err)
	assert.Equal(t, len(dockerServices), len(service.GetDependencies()))
	assert.Equal(t, service.IsRunning(), true)
	assert.Equal(t, service.IsStopped(), false)
}

func TestStatusStoped(t *testing.T) {
	service := &Service{
		Name: "TestStatusStoped",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
		},
	}
	err := service.Stop()
	assert.Nil(t, err)
	assert.Equal(t, service.IsRunning(), false)
	assert.Equal(t, service.IsStopped(), true)
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
