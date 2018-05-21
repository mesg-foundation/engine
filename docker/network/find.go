package network

import (
	godocker "github.com/fsouza/go-dockerclient"
	"github.com/mesg-foundation/core/docker"
)

// Find finds a Docker Network by a namespace
func Find(namespace string) (network *godocker.Network, err error) {
	client, err := docker.Client()
	if err != nil {
		return
	}
	networks, err := client.FilteredListNetworks(godocker.NetworkFilterOpts{
		"scope": {"swarm": true},
		"label": {"com.docker.stack.namespace=" + namespace: true},
	})
	if len(networks) > 0 {
		network = &networks[0]
	}
	return
}
