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
func (service *Service) Status() (status StatusType, err error) {
	status = STOPPED
	allRunning := true
	for dependency := range service.GetDependencies() {
		running, err := docker.IsServiceRunning([]string{service.Name, dependency})
		if err != nil {
			return status, err
		}
		if running == true {
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
func (service *Service) IsRunning() (result bool, err error) {
	status, err := service.Status()
	if err != nil {
		return
	}
	result = status == RUNNING
	return
}

// IsPartiallyRunning returns true if the service is running, false otherwise
func (service *Service) IsPartiallyRunning() (result bool, err error) {
	status, err := service.Status()
	if err != nil {
		return
	}
	result = status == PARTIAL
	return
}

// IsStopped returns true if the service is stopped, false otherwise
func (service *Service) IsStopped() (result bool, err error) {
	status, err := service.Status()
	if err != nil {
		return
	}
	result = status == STOPPED
	return
}
