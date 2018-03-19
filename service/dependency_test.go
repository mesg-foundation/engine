package service

import (
	"strings"
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/stvp/assert"
)

func TestExtractPortEmpty(t *testing.T) {
	ports := extractPorts(&Dependency{})
	assert.Equal(t, len(ports), 0)
}

func TestExtractPorts(t *testing.T) {
	ports := extractPorts(&Dependency{
		Ports: []string{
			"80",
			"3000:8080",
		},
	})
	assert.Equal(t, len(ports), 2)
	assert.Equal(t, ports[0].PublishMode, swarm.PortConfigPublishModeIngress)
	assert.Equal(t, ports[0].Protocol, swarm.PortConfigProtocolTCP)
	assert.Equal(t, ports[0].TargetPort, uint32(80))
	assert.Equal(t, ports[0].PublishedPort, uint32(80))
	assert.Equal(t, ports[1].TargetPort, uint32(8080))
	assert.Equal(t, ports[1].PublishedPort, uint32(3000))
}

func TestStartDependency(t *testing.T) {
	namespace := strings.Join([]string{NAMESPACE, "TestStartDependency"}, "_")
	name := "test"
	dependency := Dependency{Image: "nginx"}
	err := dependency.Start(name, namespace)
	assert.Nil(t, err)
	assert.Equal(t, dependency.IsRunning(namespace, name), true)
	assert.Equal(t, dependency.IsStopped(namespace, name), false)
	dependency.Stop(name, namespace)
}

func TestStopDependency(t *testing.T) {
	namespace := strings.Join([]string{NAMESPACE, "TestStopDependency"}, "_")
	name := "test"
	dependency := Dependency{Image: "nginx"}
	dependency.Start(name, namespace)
	err := dependency.Stop(name, namespace)
	assert.Nil(t, err)
	assert.Equal(t, dependency.IsStopped(namespace, name), true)
	assert.Equal(t, dependency.IsRunning(namespace, name), false)
}
