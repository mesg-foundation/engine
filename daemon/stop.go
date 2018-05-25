package daemon

import (
	"github.com/mesg-foundation/core/docker"
)

// Stop the daemon docker
func Stop() (err error) {
	stopped, err := IsStopped()
	if err != nil || stopped == true {
		return
	}
	err = docker.StopService(namespace())
	if err != nil {
		return
	}
	err = docker.DeleteNetwork(namespaceNetwork())
	if err != nil {
		return
	}
	return
}
