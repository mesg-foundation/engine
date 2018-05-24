package daemon

import (
	"github.com/mesg-foundation/core/docker"
)

// Stop the daemon docker
func Stop() (err error) {
	err = docker.StopService(Namespace())
	if err != nil {
		return
	}
	err = docker.DeleteNetwork(NamespaceNetwork())
	if err != nil {
		return
	}
	return
}
