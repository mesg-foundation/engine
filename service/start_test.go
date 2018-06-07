package service

import (
	"testing"
	"time"

	"github.com/mesg-foundation/core/container"
	"github.com/stvp/assert"
)

func TestStartService(t *testing.T) {
	service := &Service{
		Name: "TestStartService",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
		},
	}
	dockerServices, err := service.Start()
	assert.Nil(t, err)
	assert.Equal(t, len(service.GetDependencies()), len(dockerServices))
	assert.Equal(t, service.IsRunning(), true)
	service.Stop()
}

func TestStartWith2Dependencies(t *testing.T) {
	service := &Service{
		Name: "TestStartWith2Dependencies",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx:latest",
			},
			"test2": &Dependency{
				Image: "alpine:latest",
			},
		},
	}
	servicesID, err := service.Start()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(servicesID))
	container1, _ := container.FindContainer([]string{"TestStartWith2Dependencies", "test"})
	container2, _ := container.FindContainer([]string{"TestStartWith2Dependencies", "test2"})
	assert.Equal(t, "nginx:latest", container1.Config.Image)
	assert.Equal(t, "alpine:latest", container2.Config.Image)
	service.Stop()
}

func TestStartAgainService(t *testing.T) {
	service := &Service{
		Name: "TestStartAgainService",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
		},
	}
	service.Start()
	dockerServices, err := service.Start()
	assert.Nil(t, err)
	assert.Equal(t, len(dockerServices), 0) // 0 because already started so no new one to start
	assert.Equal(t, service.IsRunning(), true)
	service.Stop()
}

func TestPartiallyRunningService(t *testing.T) {
	service := &Service{
		Name: "TestPartiallyRunningService",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
			"test2": &Dependency{
				Image: "nginx",
			},
		},
	}
	service.Start()
	service.GetDependencies()["test"].Stop(service.namespace(), "test")
	assert.Equal(t, service.IsPartiallyRunning(), true)
	dockerServices, err := service.Start()
	assert.Nil(t, err)
	assert.Equal(t, len(dockerServices), len(service.GetDependencies()))
	assert.Equal(t, service.IsRunning(), true)
	service.Stop()
}

func TestStartDependency(t *testing.T) {
	namespace := container.Namespace([]string{"TestStartDependency"})
	name := "test"
	dependency := Dependency{Image: "nginx"}
	network, err := container.CreateNetwork([]string{namespace})
	serviceID, err := dependency.Start(&Service{}, dependencyDetails{
		namespace:      namespace,
		dependencyName: name,
		serviceName:    "TestStartDependency",
	}, network)
	assert.Nil(t, err)
	assert.NotEqual(t, "", serviceID)
	assert.Equal(t, dependency.IsRunning(namespace, name), true)
	assert.Equal(t, dependency.IsStopped(namespace, name), false)
	dependency.Stop(namespace, name)
	container.DeleteNetwork([]string{namespace})
}

func TestNetworkCreated(t *testing.T) {
	service := &Service{
		Name: "TestNetworkCreated",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
		},
	}
	service.Start()
	network, err := container.FindNetwork([]string{service.namespace()})
	assert.Nil(t, err)
	assert.NotEqual(t, "", network.ID)
	service.Stop()
}

// Test for https://github.com/mesg-foundation/core/issues/88
func TestStartStopStart(t *testing.T) {
	service := &Service{
		Name: "TestStartStopStart",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
		},
	}
	service.Start()
	service.Stop()
	time.Sleep(10 * time.Second)
	dockerServices, err := service.Start()
	assert.Nil(t, err)
	assert.Equal(t, len(dockerServices), 1)
	assert.Equal(t, service.IsRunning(), true)
	service.Stop()
}
