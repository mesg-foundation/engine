package service

import (
	"github.com/mesg-foundation/core/docker"
)

// StatusType of the service
type StatusType uint

// status for services
const (
	STOPPED StatusType = 1
	RUNNING StatusType = 2
	PARTIAL StatusType = 3
)

// Status return the status of the docker service for this service
func Status(service *Service) (status StatusType) {
	status = STOPPED
	allRunning := true
	for dependency := range service.GetDependencies() {
		if docker.IsServiceRunning([]string{service.Name, dependency}) {
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
	running = Status(service) == RUNNING
	return
}

// IsPartiallyRunning returns true if the service is running, false otherwise
func (service *Service) IsPartiallyRunning() (running bool) {
	running = Status(service) == PARTIAL
	return
}

// IsStopped returns true if the service is stopped, false otherwise
func (service *Service) IsStopped() (running bool) {
	running = Status(service) == STOPPED
	return
}
