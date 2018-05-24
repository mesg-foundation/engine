package service

import (
	"github.com/docker/docker/api/types/swarm"
	"github.com/mesg-foundation/core/docker"
)

// List all the mesg docker services
func List() (services []swarm.Service, err error) {
	services, err = docker.ListServices(dockerLabelServiceKey)
	if err != nil {
		return
	}
	return
}
