package daemon

import (
	"github.com/mesg-foundation/core/docker"
)

// Stop the daemon docker
func Stop() (err error) {
	err = docker.StopService([]string{name})
	if err != nil {
		return
	}
	err = docker.DeleteNetwork([]string{sharedNetwork})
	if err != nil {
		return
	}
	return
}
