package service

import (
	"fmt"
	"testing"

	"github.com/mesg-foundation/core/docker"
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
	fmt.Println(err)
	assert.Nil(t, err)
	assert.Equal(t, len(dockerServices), len(service.GetDependencies()))
	assert.Equal(t, service.IsRunning(), true)
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
	docker.StopService([]string{service.Name, "test"})
	assert.Equal(t, service.IsPartiallyRunning(), true)
	dockerServices, err := service.Start()
	assert.Nil(t, err)
	assert.Equal(t, len(dockerServices), len(service.GetDependencies()))
	assert.Equal(t, service.IsRunning(), true)
	service.Stop()
}

func TestStartDependency(t *testing.T) {
	c := dockerConfig{
		service: &Service{
			Name: "TestStartDependency",
		},
		dependency: &Dependency{
			Image: "nginx",
		},
		name: "test",
	}
	namespaces := []string{c.service.Name, c.name}
	dockerService, err := dockerStart(c)
	assert.Nil(t, err)
	assert.NotNil(t, dockerService)
	assert.Equal(t, docker.IsServiceRunning(namespaces), true)
	assert.Equal(t, docker.IsServiceStopped(namespaces), false)
	docker.StopService(namespaces)
}

// func TestNetworkCreated(t *testing.T) {
// 	service := &Service{
// 		Name: "TestNetworkCreated",
// 		Dependencies: map[string]*Dependency{
// 			"test": &Dependency{
// 				Image: "nginx",
// 			},
// 		},
// 	}
// 	service.Start(testDaemonIP, testSharedNetwork)
// 	network, err := findNetwork(service.namespace())
// 	assert.Nil(t, err)
// 	assert.NotNil(t, network)
// 	service.Stop()
// }

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
	dockerServices, err := service.Start()
	assert.Nil(t, err)
	assert.Equal(t, len(dockerServices), 1)
	assert.Equal(t, service.IsRunning(), true)
	service.Stop()
}
