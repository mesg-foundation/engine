package service

import (
	"github.com/mesg-foundation/core/docker"
)

// Status return the Docker service status of a service
func Status(service service) (status docker.StatusType) {
	status = docker.STOPPED
	allRunning := true
	keys := service.GetDependenciesKeys()
	for _, name := range keys {
		if docker.IsRunning(service.Namespace(), name) {
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
func IsRunning(service service) (running bool) {
	running = Status(service) == docker.RUNNING
	return
}

// IsPartiallyRunning returns true if the service is running, false otherwise
func IsPartiallyRunning(service service) (running bool) {
	running = Status(service) == docker.PARTIAL
	return
}

// IsStopped returns true if the service is stopped, false otherwise
func IsStopped(service service) (running bool) {
	running = Status(service) == docker.STOPPED
	return
}
