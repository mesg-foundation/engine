package service

import (
	"github.com/mesg-foundation/core/docker"
)

// Stop a service
func (service *Service) Stop() (err error) {
	stopped, err := service.IsStopped()
	if err != nil || stopped == true {
		return
	}
	for dependency := range service.Dependencies {
		err = docker.StopService([]string{service.Name, dependency})
		if err != nil {
			return
		}
	}
	err = docker.DeleteNetwork([]string{service.Name})
	if err != nil {
		return
	}
	return
}
