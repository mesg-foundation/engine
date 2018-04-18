package service

import (
	"os"
	"testing"

	"github.com/stvp/assert"
)

func TestStatusRunning(t *testing.T) {
	// TODO remove and make CI works
	if os.Getenv("CI") == "true" {
		return
	}
	service := &Service{
		Name: "TestStatusRunning",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
		},
	}
	dockerServices, err := service.Start()
	assert.Nil(t, err)
	assert.Equal(t, len(dockerServices), len(service.Dependencies))
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
