package service

import (
	"os"
	"strings"
	"testing"

	"github.com/stvp/assert"
)

func TestStartService(t *testing.T) {
	// TODO remove and make CI works
	if os.Getenv("CI") == "true" {
		return
	}
	service := &Service{
		Name: "TestStartService",
		Dependencies: map[string]Dependency{
			"test": Dependency{
				Image: "nginx",
			},
		},
	}
	dockerServices, err := service.Start()
	assert.Nil(t, err)
	assert.Equal(t, len(dockerServices), len(service.Dependencies))
	assert.Equal(t, service.IsRunning(), true)
	service.Stop()
}

func TestStartAgainService(t *testing.T) {
	// TODO remove and make CI works
	if os.Getenv("CI") == "true" {
		return
	}
	service := &Service{
		Name: "TestStartAgainService",
		Dependencies: map[string]Dependency{
			"test": Dependency{
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
	// TODO remove and make CI works
	if os.Getenv("CI") == "true" {
		return
	}
	service := &Service{
		Name: "TestPartiallyRunningService",
		Dependencies: map[string]Dependency{
			"test": Dependency{
				Image: "nginx",
			},
			"test2": Dependency{
				Image: "nginx",
			},
		},
	}
	service.Start()
	service.Dependencies["test"].Stop(service.namespace(), "test")
	assert.Equal(t, service.IsPartiallyRunning(), true)
	dockerServices, err := service.Start()
	assert.Nil(t, err)
	assert.Equal(t, len(dockerServices), len(service.Dependencies))
	assert.Equal(t, service.IsRunning(), true)
	service.Stop()
}

func TestStartDependency(t *testing.T) {
	// TODO remove and make CI works
	if os.Getenv("CI") == "true" {
		return
	}
	namespace := strings.Join([]string{NAMESPACE, "TestStartDependency"}, "_")
	name := "test"
	dependency := Dependency{Image: "nginx"}
	network, err := createNetwork(namespace)
	dockerService, err := dependency.Start(namespace, name, network)
	assert.Nil(t, err)
	assert.NotNil(t, dockerService)
	assert.Equal(t, dependency.IsRunning(namespace, name), true)
	assert.Equal(t, dependency.IsStopped(namespace, name), false)
	dependency.Stop(namespace, name)
}
