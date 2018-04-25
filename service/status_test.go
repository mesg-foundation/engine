package service

import (
	"os"
	"testing"

	"github.com/mesg-foundation/core/types"
	"github.com/stvp/assert"
)

func TestStatusRunning(t *testing.T) {
	// TODO remove and make CI works
	if os.Getenv("CI") == "true" {
		return
	}
	service := &Service{
		Name: "TestStatusRunning",
		Dependencies: map[string]*types.ProtoDependency{
			"test": &types.ProtoDependency{
				Image: "nginx",
			},
		},
	}
	dockerServices, err := service.Start()
	assert.Nil(t, err)
	assert.Equal(t, len(dockerServices), len(service.GetDependencies()))
	assert.Equal(t, service.IsRunning(), true)
	assert.Equal(t, service.IsStopped(), false)
	service.Stop()
}

func TestStatusStoped(t *testing.T) {
	if os.Getenv("CI") == "true" {
		return
	}
	service := &Service{
		Name: "TestStatusStoped",
		Dependencies: map[string]*types.ProtoDependency{
			"test": &types.ProtoDependency{
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
	if os.Getenv("CI") == "true" {
		return
	}
	service := &Service{
		Name: "TestList",
		Dependencies: map[string]*types.ProtoDependency{
			"test": &types.ProtoDependency{
				Image: "nginx",
			},
		},
	}
	service.Start()
	list, err := List()
	assert.Nil(t, err)
	assert.Equal(t, len(list), 1)
	assert.Equal(t, list[0], service.Name)
	service.Stop()
}

func TestListMultipleDependencies(t *testing.T) {
	if os.Getenv("CI") == "true" {
		return
	}
	service := &Service{
		Name: "TestListMultipleDependencies",
		Dependencies: map[string]*types.ProtoDependency{
			"test": &types.ProtoDependency{
				Image: "nginx",
			},
			"test2": &types.ProtoDependency{
				Image: "nginx",
			},
		},
	}
	service.Start()
	list, err := List()
	assert.Nil(t, err)
	assert.Equal(t, len(list), 1)
	assert.Equal(t, list[0], service.Name)
	service.Stop()
}
