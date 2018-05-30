package daemon

import (
	"github.com/mesg-foundation/core/container"

	"github.com/docker/docker/api/types/swarm"
	"github.com/fsouza/go-dockerclient"
)

func getContainer() (*docker.APIContainers, error) {
	return container.FindContainer([]string{name})
}

func getService() (*swarm.Service, error) {
	return container.FindService([]string{name})
}
