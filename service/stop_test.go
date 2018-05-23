package service

import (
	"testing"

	"github.com/mesg-foundation/core/docker"
	"github.com/stvp/assert"
)

func TestStopRunningService(t *testing.T) {
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
	assert.Equal(t, service.IsStopped(), true)
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
	assert.Equal(t, service.IsStopped(), true)
}

func TestStopDependency(t *testing.T) {
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
	assert.Equal(t, docker.IsServiceStopped(namespaces), true)
	assert.Equal(t, docker.IsServiceRunning(namespaces), false)
}

func TestNetworkDeleted(t *testing.T) {
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
	network, err := docker.FindNetwork([]string{service.Name})
	assert.Nil(t, err)
	assert.Nil(t, network)
}
