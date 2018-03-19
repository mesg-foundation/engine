package service

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types/swarm"
	docker "github.com/fsouza/go-dockerclient"
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

func dependencyStatus(dependency *Dependency, namespace string, dependencyName string) (status StatusType) {
	ctx := context.Background()
	dockerServices, err := dockerCli.ListServices(docker.ListServicesOptions{
		Filters: map[string][]string{
			"name": []string{strings.Join([]string{namespace, dependencyName}, "_")},
		},
		Context: ctx,
	})
	dockerService := dockerServiceMatch(dockerServices, namespace, dependencyName)
	status = STOPPED
	if err == nil && dockerService.ID != "" {
		status = RUNNING
	}
	return
}

func serviceStatus(service *Service) (status StatusType) {
	status = STOPPED
	allRunning := true
	for name, dependency := range service.Dependencies {
		if dependency.IsRunning(service.namespace(), name) {
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
	status := serviceStatus(service)
	running = status == RUNNING
	return
}

// IsPartiallyRunning returns true if the service is running, false otherwise
func (service *Service) IsPartiallyRunning() (running bool) {
	status := serviceStatus(service)
	running = status == PARTIAL
	return
}

// IsStopped returns true if the service is stopped, false otherwise
func (service *Service) IsStopped() (running bool) {
	status := serviceStatus(service)
	running = status == STOPPED
	return
}

// IsRunning returns true if the dependency is running, false otherwise
func (dependency *Dependency) IsRunning(namespace string, name string) (running bool) {
	status := dependencyStatus(dependency, namespace, name)
	running = status == RUNNING
	return
}

// IsStopped returns true if the dependency is stopped, false otherwise
func (dependency *Dependency) IsStopped(namespace string, name string) (running bool) {
	status := dependencyStatus(dependency, namespace, name)
	running = status == STOPPED
	return
}
