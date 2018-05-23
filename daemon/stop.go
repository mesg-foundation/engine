package daemon

import (
	"github.com/mesg-foundation/core/docker"
)

func Stop() (err error) {
	service, err := docker.FindService([]string{name})
	if err != nil {
		return
	}
	if service != nil {
		err = docker.StopService([]string{name})
		if err != nil {
			return
		}
	}

	network, err := docker.FindNetwork(sharedNetwork)
	if err != nil {
		return
	}
	if network != nil {
		err = docker.DeleteNetwork(sharedNetwork)
		if err != nil {
			return
		}
	}

	return
}
