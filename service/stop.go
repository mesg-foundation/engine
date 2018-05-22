package service

import (
	"github.com/mesg-foundation/core/docker"
	dockerService "github.com/mesg-foundation/core/docker/service"
)

// Stop a service
func (service *Service) Stop() (err error) {
	if dockerService.IsStopped(service) {
		return
	}
	for _, name := range service.GetDependenciesKeys() {
		err = docker.Stop(service.Namespace(), name)
		if err != nil {
			break
		}
	}
	// TODO: docker shared network: remove the specific docker network for this docker service
	// if err == nil { // didnt exit the loop
	// 	err = deleteNetwork(service.namespace())
	// }
	return
}
