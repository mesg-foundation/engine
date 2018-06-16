package service

import (
	"testing"
	"time"

	"github.com/mesg-foundation/core/container"
	"github.com/stvp/assert"
)

func TestExtractPortEmpty(t *testing.T) {
	dep := Dependency{}
	ports := dep.extractPorts()
	assert.Equal(t, len(ports), 0)
}

func TestExtractPorts(t *testing.T) {
	dep := &Dependency{
		Ports: []string{
			"80",
			"3000:8080",
		},
	}
	ports := dep.extractPorts()
	assert.Equal(t, len(ports), 2)
	assert.Equal(t, ports[0].Target, uint32(80))
	assert.Equal(t, ports[0].Published, uint32(80))
	assert.Equal(t, ports[1].Target, uint32(8080))
	assert.Equal(t, ports[1].Published, uint32(3000))
}

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
	deps := service.DependenciesFromService()
	container1, _ := container.FindContainer(deps[0].namespace())
	container2, _ := container.FindContainer(deps[1].namespace())
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
	service.DependenciesFromService()[0].Stop()
	assert.Equal(t, service.IsPartiallyRunning(), true)
	dockerServices, err := service.Start()
	assert.Nil(t, err)
	assert.Equal(t, len(dockerServices), len(service.GetDependencies()))
	assert.Equal(t, service.IsRunning(), true)
	service.Stop()
}

func TestStartDependency(t *testing.T) {
	service := &Service{
		Name: "TestStartDependency",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
		},
	}
	networkID, err := container.CreateNetwork(service.namespace())
	dep := service.DependenciesFromService()[0]
	serviceID, err := dep.Start(networkID)
	assert.Nil(t, err)
	assert.NotEqual(t, "", serviceID)
	assert.Equal(t, dep.IsRunning(), true)
	assert.Equal(t, dep.IsStopped(), false)
	dep.Stop()
	container.DeleteNetwork(service.namespace())
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
	network, err := container.FindNetwork(service.namespace())
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
