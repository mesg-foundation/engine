package service

import (
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/mesg-foundation/core/daemon"
	"github.com/mesg-foundation/core/docker"
	"github.com/stvp/assert"
)

func startTestService(name string, dependency string) (service *Service, swarmService []*swarm.Service, err error) {
	daemon.Start()
	service = &Service{
		Name: name,
		Dependencies: map[string]*Dependency{
			dependency: &Dependency{
				Image: "nginx",
			},
		},
	}
	swarmService, err = service.Start()
	return
}

func TestPortsEmpty(t *testing.T) {
	c := dockerConfig{
		dependency: &Dependency{},
	}
	ports := c.dockerPorts()
	assert.Equal(t, len(ports), 0)
}

func TestPorts(t *testing.T) {
	c := dockerConfig{
		dependency: &Dependency{
			Ports: []string{
				"80",
				"3000:8080",
			},
		},
	}
	ports := c.dockerPorts()
	assert.Equal(t, len(ports), 2)
	assert.Equal(t, ports[0].Target, uint32(80))
	assert.Equal(t, ports[0].Published, uint32(80))
	assert.Equal(t, ports[1].Target, uint32(8080))
	assert.Equal(t, ports[1].Published, uint32(3000))
}

func TestStartService(t *testing.T) {
	daemon.Start()
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
	assert.Equal(t, len(dockerServices), len(service.GetDependencies()))
	running, err := service.IsRunning()
	assert.Nil(t, err)
	assert.Equal(t, running, true)
	service.Stop()
}

func TestStartAgainService(t *testing.T) {
	daemon.Start()
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
	running, err := service.IsRunning()
	assert.Nil(t, err)
	assert.Equal(t, running, true)
	service.Stop()
}

func TestPartiallyRunningService(t *testing.T) {
	daemon.Start()
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
	_, err := service.Start()
	assert.Nil(t, err)
	docker.StopService([]string{service.Name, "test"})
	partial, err := service.IsPartiallyRunning()
	assert.Nil(t, err)
	assert.Equal(t, partial, true)
	dockerServices, err := service.Start()
	assert.Nil(t, err)
	assert.Equal(t, len(dockerServices), len(service.GetDependencies()))
	running, err := service.IsRunning()
	assert.Nil(t, err)
	assert.Equal(t, running, true)
	service.Stop()
}

func TestStartDependency(t *testing.T) {
	daemon.Start()
	c := dockerConfig{
		service: &Service{
			Name: "TestStartDependency",
		},
		dependency: &Dependency{
			Image: "nginx",
		},
		name:      "test",
		networkID: "host",
	}
	namespaces := []string{c.service.Name, c.name}
	dockerService, err := startDocker(c)
	assert.Nil(t, err)
	assert.NotNil(t, dockerService)
	running, err := docker.IsServiceRunning(namespaces)
	assert.Nil(t, err)
	assert.Equal(t, running, true)
	stopped, err := docker.IsServiceStopped(namespaces)
	assert.Nil(t, err)
	assert.Equal(t, stopped, false)
	docker.StopService(namespaces)
}

func TestNetworkCreated(t *testing.T) {
	daemon.Start()
	service := &Service{
		Name: "TestNetworkCreated",
		Dependencies: map[string]*Dependency{
			"test": &Dependency{
				Image: "nginx",
			},
		},
	}
	service.Start()
	network, err := docker.FindNetwork([]string{service.Name})
	assert.Nil(t, err)
	assert.NotNil(t, network)
	service.Stop()
}

// Test for https://github.com/mesg-foundation/core/issues/88
func TestStartStopStart(t *testing.T) {
	daemon.Start()
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
	dockerServices, err := service.Start()
	assert.Nil(t, err)
	assert.Equal(t, len(dockerServices), 1)
	running, err := service.IsRunning()
	assert.Nil(t, err)
	assert.Equal(t, running, true)
	service.Stop()
}
