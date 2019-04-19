package dockermanager

import "github.com/mesg-foundation/core/container"

// DockerManager is responsible for managing MESG Service's Docker Containers
// in Docker Service environment.
type DockerManager struct {
	c container.Container
}

// New returns a new Docker Manager.
func New(c container.Container) *DockerManager {
	return &DockerManager{
		c: c,
	}
}
