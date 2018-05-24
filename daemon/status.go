package daemon

import (
	"time"

	"github.com/mesg-foundation/core/docker"
)

// IsRunning returns true if the daemon container is running
func IsRunning() (running bool, err error) {
	status, err := docker.ContainerStatus(Namespace())
	if err != nil {
		return
	}
	running = status == docker.RUNNING
	return
}

// WaitForRunning wait for the Daemon container to run
func WaitForRunning() (wait chan error) {
	return docker.WaitContainerStatus(Namespace(), docker.RUNNING, 10*time.Minute)
}

// WaitForStopped wait for the Daemon container to stop
func WaitForStopped() (wait chan error) {
	return docker.WaitContainerStatus(Namespace(), docker.STOPPED, 10*time.Minute)
}
