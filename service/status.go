package service

import (
	"strings"

	"github.com/docker/docker/api/types/swarm"
)

// StatusType of the service
type StatusType uint

// status for services
const (
	STOPPED StatusType = 0
	RUNNING StatusType = 1
	PARTIAL StatusType = 2
)

func dockerServiceMatch(dockerServices []swarm.Service, namespace string, name string) (dockerService swarm.Service) {
	for _, service := range dockerServices {
		if service.Spec.Annotations.Name == strings.Join([]string{namespace, name}, "_") {
			dockerService = service
			break
		}
	}
	return
}

func serviceStatus(service *Service) (status StatusType) {
	status = STOPPED
	allRunning := true
	for name, dependency := range service.Dependencies {
		if dependency.IsRunning(service.Namespace(), name) {
			status = RUNNING
		} else {
			allRunning = false
		}
	}
	if status == RUNNING && !allRunning {
		status = PARTIAL
	}
	return
}

// IsRunning returns true if the service is running, false otherwise
func (service *Service) IsRunning() (running bool) {
	running = serviceStatus(service) == RUNNING
	return
}

// IsPartiallyRunning returns true if the service is running, false otherwise
func (service *Service) IsPartiallyRunning() (running bool) {
	running = serviceStatus(service) == PARTIAL
	return
}

// IsStopped returns true if the service is stopped, false otherwise
func (service *Service) IsStopped() (running bool) {
	running = serviceStatus(service) == STOPPED
	return
}

// IsRunning returns true if the dependency is running, false otherwise
func (dependency Dependency) IsRunning(namespace string, name string) (running bool) {
	running = dependencyStatus(namespace, name) == RUNNING
	return
}

// IsStopped returns true if the dependency is stopped, false otherwise
func (dependency Dependency) IsStopped(namespace string, name string) (running bool) {
	running = dependencyStatus(namespace, name) == STOPPED
	return
}
