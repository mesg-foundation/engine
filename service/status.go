package service

import (
	"errors"
	"time"

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

// WaitStatus waits until the service to have the given status until it reach the timeout
func (service *Service) WaitStatus(status StatusType, timeout time.Duration) (wait chan error) {
	start := time.Now()
	wait = make(chan error, 1)
	go func() {
		for {
			currentStatus, err := service.Status()
			if err != nil {
				wait <- err
				return
			}
			if currentStatus == status {
				close(wait)
				return
			}
			diff := time.Now().Sub(start)
			if diff.Nanoseconds() >= int64(timeout) {
				wait <- errors.New("Wait too long for the service, timeout reached")
				return
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
	return
}
