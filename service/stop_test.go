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
	namespace := container.Namespace([]string{"TestStopDependency"})
	name := "test"
	dependency := Dependency{Image: "nginx"}

	network, err := container.CreateNetwork([]string{namespace})
	dependency.Start(&Service{}, dependencyDetails{
		namespace:      namespace,
		dependencyName: name,
		serviceName:    "TestStopDependency",
	}, network)
	err = dependency.Stop(namespace, name)
	assert.Nil(t, err)
	assert.Equal(t, dependency.IsStopped(namespace, name), true)
	assert.Equal(t, dependency.IsRunning(namespace, name), false)
	container.DeleteNetwork([]string{namespace})
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
	_, err := container.FindNetwork([]string{service.namespace()})
	assert.NotNil(t, err)
}
