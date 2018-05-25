package service

import (
	"testing"

	"github.com/mesg-foundation/core/daemon"
	"github.com/stvp/assert"
)

func TestStatusRunning(t *testing.T) {
	daemon.Start()
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
	status, err := service.Status()
	assert.Nil(t, err)
	assert.Equal(t, RUNNING, status)
	running, err := service.IsRunning()
	assert.Nil(t, err)
	assert.Equal(t, true, running)
	partial, err := service.IsPartiallyRunning()
	assert.Nil(t, err)
	assert.Equal(t, false, partial)
	stopped, err := service.IsStopped()
	assert.Nil(t, err)
	assert.Equal(t, false, stopped)
	service.Stop()
}

func TestStatusStopped(t *testing.T) {
	daemon.Start()
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
	status, err := service.Status()
	assert.Nil(t, err)
	assert.Equal(t, STOPPED, status)
	running, err := service.IsRunning()
	assert.Nil(t, err)
	assert.Equal(t, false, running)
	partial, err := service.IsPartiallyRunning()
	assert.Nil(t, err)
	assert.Equal(t, false, partial)
	stopped, err := service.IsStopped()
	assert.Nil(t, err)
	assert.Equal(t, true, stopped)
}
