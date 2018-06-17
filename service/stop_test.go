package service

import (
	"testing"
	"time"

	"github.com/mesg-foundation/core/container"
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
	service := &Service{
		Name: "TestStartDependency",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
		},
	}
	networkID, err := container.CreateNetwork(service.namespace())
	defer container.DeleteNetwork(service.namespace())
	dep := service.DependenciesFromService()[0]
	dep.Start(networkID)
	err = dep.Stop()
	assert.Nil(t, err)
	assert.Equal(t, dep.IsStopped(), true)
	assert.Equal(t, dep.IsRunning(), false)
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
	time.Sleep(5 * time.Second)
	_, err := container.FindNetwork(service.namespace())
	assert.NotNil(t, err)
}
