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
	for _, dependency := range service.DependenciesFromService() {
		if dependency.IsRunning() {
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
func (dependency *DependencyFromService) IsRunning() (running bool) {
	running = dependency.Status() == RUNNING
	return
}

// IsStopped returns true if the dependency is stopped, false otherwise
func (dependency *DependencyFromService) IsStopped() (running bool) {
	running = dependency.Status() == STOPPED
	return
}

// Status returns the status of this dependency's container
func (dependency *DependencyFromService) Status() (status StatusType) {
	dockerStatus, err := container.ServiceStatus(dependency.namespace())
	if err != nil {
		panic(err) //TODO: that's VERY ugly
	}
	status = STOPPED
	if dockerStatus == container.RUNNING {
		status = RUNNING
	}
	return
}

// ListRunning all the running services
// TODO: should move to another file
func ListRunning() (res []string, err error) {
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
