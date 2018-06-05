package service

import (
	"github.com/mesg-foundation/core/container"
)

// StatusType of the service
type StatusType uint

// status for services
const (
	STOPPED StatusType = 0
	RUNNING StatusType = 1
	PARTIAL StatusType = 2
)

func serviceStatus(service *Service) (status StatusType) {
	status = STOPPED
	allRunning := true
	for name, dependency := range service.GetDependencies() {
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
func (dependency *Dependency) IsRunning(namespace string, name string) (running bool) {
	running = dependencyStatus(namespace, name) == RUNNING
	return
}

// IsStopped returns true if the dependency is stopped, false otherwise
func (dependency *Dependency) IsStopped(namespace string, name string) (running bool) {
	running = dependencyStatus(namespace, name) == STOPPED
	return
}

// List all the running services hashes
func List() (res []string, err error) {
	services, err := container.ListServices("mesg.hash")
	mapRes := make(map[string]uint)
	for _, service := range services {
		serviceName := service.Spec.Annotations.Labels["mesg.hash"]
		mapRes[serviceName]++
	}
	res = make([]string, 0, len(mapRes))
	for k := range mapRes {
		res = append(res, k)
	}
	return
}
