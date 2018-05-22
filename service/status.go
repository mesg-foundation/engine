package service

import (
	"github.com/mesg-foundation/core/docker"
)

// Status return the status of the docker service for this service
func Status(service *Service) (status docker.StatusType) {
	status = docker.STOPPED
	allRunning := true
	for dependency := range service.GetDependencies() {
		if docker.IsRunning([]string{service.Name, dependency}) {
			status = docker.RUNNING
		} else {
			allRunning = false
		}
	}
	if status == docker.RUNNING && !allRunning {
		status = docker.PARTIAL
	}
	return
}

// IsRunning returns true if the service is running, false otherwise
func (service *Service) IsRunning() (running bool) {
	running = Status(service) == docker.RUNNING
	return
}

// IsPartiallyRunning returns true if the service is running, false otherwise
func (service *Service) IsPartiallyRunning() (running bool) {
	running = Status(service) == docker.PARTIAL
	return
}

// IsStopped returns true if the service is stopped, false otherwise
func (service *Service) IsStopped() (running bool) {
	running = Status(service) == docker.STOPPED
	return
}
