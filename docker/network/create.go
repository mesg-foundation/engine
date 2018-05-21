package network

import (
	"strings"

	godocker "github.com/fsouza/go-dockerclient"
	"github.com/mesg-foundation/core/docker"
)

// Create creates a Docker Network with a namespace
func Create(namespace string) (network *godocker.Network, err error) {
	client, err := docker.Client()
	if err != nil {
		return
	}
	network, err = Find(namespace)
	if network != nil || err != nil {
		return
	}
	network, err = client.CreateNetwork(godocker.CreateNetworkOptions{
		Name:   strings.Join([]string{namespace, "Network"}, "-"),
		Driver: "overlay",
		Labels: map[string]string{
			"com.docker.stack.namespace": namespace,
		},
	})
	return
}
