package daemon

import "github.com/mesg-foundation/core/docker"

// IsRunning returns true if the daemon container is running
func IsRunning() (running bool, err error) {
	container, err := docker.FindContainer(image)
	if err != nil {
		return
	}
	running = container != nil
	return
}
