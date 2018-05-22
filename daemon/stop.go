package daemon

import (
	"github.com/mesg-foundation/core/docker"
)

func Stop() (err error) {
	service, err := Service()
	if err != nil {
		return
	}
	if service != nil {
		err = docker.Stop([]string{name})
		if err != nil {
			return
		}
	}

	network, err := Network()
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
