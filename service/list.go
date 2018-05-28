package service

import (
	"github.com/docker/docker/api/types/swarm"
	"github.com/mesg-foundation/core/docker"
)

// List all the mesg docker services
// TODO: this function should return the MESG Service not the docker service (can have multiple docker service per mesg service)
func List() (services []swarm.Service, err error) {
	services, err = docker.ListServices(dockerLabelServiceKey)
	if err != nil {
		return
	}
	return
}
