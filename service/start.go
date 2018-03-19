package service

import "github.com/docker/docker/api/types/swarm"

// Start a service
func (service *Service) Start() (dockerServices []*swarm.Service, err error) {
	if service.IsRunning() {
		return
	}
	// If there is one but not all services running stop to restart all
	if service.IsPartiallyRunning() {
		service.Stop()
	}
	dockerServices = make([]*swarm.Service, len(service.Dependencies))
	i := 0
	for name, dependency := range service.Dependencies {
		dockerServices[i], err = dependency.Start(service.namespace(), name)
		i++
		if err != nil {
			break
		}
	}
	// Disgrasfully close the service because there is an error
	if err != nil {
		service.Stop()
	}
	return
}
