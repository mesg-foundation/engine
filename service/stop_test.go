package service

import (
	"testing"
	"time"

	"github.com/mesg-foundation/core/daemon"
	"github.com/mesg-foundation/core/docker"
	"github.com/stvp/assert"
)

func TestStopRunningService(t *testing.T) {
	daemon.Start()
	service := &Service{
		Name: "TestStopRunningService",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
		},
	}
	service.Start()
	err := service.Stop()
	assert.Nil(t, err)
	stopped, err := service.IsStopped()
	assert.Nil(t, err)
	assert.Equal(t, true, stopped)
}

func TestStopNonRunningService(t *testing.T) {
	service := &Service{
		Name: "TestStopNonRunningService",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
		},
	}
	err := service.Stop()
	assert.Nil(t, err)
	stopped, err := service.IsStopped()
	assert.Nil(t, err)
	assert.Equal(t, true, stopped)
}

func TestStopDependency(t *testing.T) {
	daemon.Start()
	c := dockerConfig{
		service: &Service{
			Name: "TestStopDependency",
		},
		dependency: &Dependency{
			Image: "nginx",
		},
		name: "test",
	}
	namespaces := []string{c.service.Name, c.name}
	startDocker(c)
	err := docker.StopService(namespaces)
	assert.Nil(t, err)
	stopped, err := docker.IsServiceStopped(namespaces)
	assert.Nil(t, err)
	assert.Equal(t, true, stopped)
}

func TestNetworkDeleted(t *testing.T) {
	daemon.Start()
	service := &Service{
		Name: "TestNetworkDeleted",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
		},
	}
	service.Start()
	service.Stop()
	<-service.WaitStatus(STOPPED, 30*time.Second)
	network, err := docker.FindNetwork([]string{service.Name})
	assert.Nil(t, err)
	assert.Nil(t, network)
}
