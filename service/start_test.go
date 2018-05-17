package service

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stvp/assert"
)

var (
	testDaemonIP      = "localhost" // TODO: should be remove when a better implementation is up
	testSharedNetwork = "ingress"   // TODO: should be remove when a better implementation is up
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
	dockerServices, err := service.Start(testDaemonIP, testSharedNetwork)
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
	service.Start(testDaemonIP, testSharedNetwork)
	dockerServices, err := service.Start(testDaemonIP, testSharedNetwork)
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
	service.Start(testDaemonIP, testSharedNetwork)
	service.GetDependencies()["test"].Stop(service.namespace(), "test")
	assert.Equal(t, service.IsPartiallyRunning(), true)
	dockerServices, err := service.Start(testDaemonIP, testSharedNetwork)
	assert.Nil(t, err)
	assert.Equal(t, len(dockerServices), len(service.GetDependencies()))
	assert.Equal(t, service.IsRunning(), true)
	service.Stop()
}

func TestStartDependency(t *testing.T) {
	namespace := strings.Join([]string{NAMESPACE, "TestStartDependency"}, "_")
	name := "test"
	dependency := Dependency{Image: "nginx"}
	dockerService, err := dependency.Start(&Service{}, dependencyDetails{
		namespace:      namespace,
		dependencyName: name,
		serviceName:    "TestStartDependency",
	}, testDaemonIP, testSharedNetwork)
	assert.Nil(t, err)
	assert.NotNil(t, dockerService)
	assert.Equal(t, dependency.IsRunning(namespace, name), true)
	assert.Equal(t, dependency.IsStopped(namespace, name), false)
	dependency.Stop(namespace, name)
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
	service.Start(testDaemonIP, testSharedNetwork)
	service.Stop()
	dockerServices, err := service.Start(testDaemonIP, testSharedNetwork)
	assert.Nil(t, err)
	assert.Equal(t, len(dockerServices), 1)
	assert.Equal(t, service.IsRunning(), true)
	service.Stop()
}
