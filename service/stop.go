package service

import (
	"github.com/mesg-foundation/core/docker"
)

// Stop a service
func (service *Service) Stop() (err error) {
	if service.IsStopped() {
		return
	}
	for dependency := range service.Dependencies {
		err = docker.StopService([]string{service.Name, dependency})
		if err != nil {
			return
		}
	}
	// TODO: docker shared network: remove the specific docker network for this docker service
	// if err == nil { // didnt exit the loop
	// 	err = deleteNetwork(service.namespace())
	// }
	return
}
