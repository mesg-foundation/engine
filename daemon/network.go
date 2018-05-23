package daemon

import (
	"errors"

	"github.com/mesg-foundation/core/docker"
)

// IP returns the IP of the daemon in the shared network
func IP() (daemonIP string, err error) {
	daemonContainer, err := docker.FindContainer(name)
	if err != nil {
		return
	}
	if daemonContainer == nil {
		err = errors.New("Daemon container not found")
		return
	}
	networkContainer := daemonContainer.Networks.Networks["mesg-shared-network"]
	if networkContainer.IPAddress == "" {
		err = errors.New("Network 'mesg-shared-network' not found")
		return
	}
	daemonIP = networkContainer.IPAddress
	return
}

// SharedNetworkID returns the shared network id
func SharedNetworkID() (networkID string, err error) {
	network, err := docker.FindNetwork(sharedNetwork)
	if err != nil {
		return
	}
	networkID = network.ID
	return
}
