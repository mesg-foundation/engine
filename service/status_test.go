package service

import (
	"testing"

	"github.com/stvp/assert"
)

func TestStatusRunning(t *testing.T) {
	testStartDaemon()
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
	assert.Equal(t, len(dockerServices), len(service.GetDependencies()))
	assert.Equal(t, true, service.IsRunning())
	assert.Equal(t, false, service.IsStopped())
	service.Stop()
}

func TestStatusStopped(t *testing.T) {
	testStartDaemon()
	service := &Service{
		Name: "TestStatusStopped",
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
