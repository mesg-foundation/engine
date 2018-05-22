package docker

import (
	"strings"

	godocker "github.com/fsouza/go-dockerclient"
)

// Create creates a Docker Network with a namespace
func NetworkCreate(namespace string) (network *godocker.Network, err error) {
	client, err := Client()
	if err != nil {
		return
	}
	network, err = NetworkFind(namespace)
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

// Delete deletes a Docker Network associated with a namespace
func NetworkDelete(namespace string) (err error) {
	client, err := Client()
	if err != nil {
		return
	}
	network, err := NetworkFind(namespace)
	if err != nil {
		return
	}
	return client.RemoveNetwork(network.ID)
}

// Find finds a Docker Network by a namespace
func NetworkFind(namespace string) (network *godocker.Network, err error) {
	client, err := Client()
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
