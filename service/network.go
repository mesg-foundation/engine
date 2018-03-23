package service

import (
	"strings"

	"github.com/fsouza/go-dockerclient"
)

func findNetwork(namespace string) (network *docker.Network, err error) {
	networks, err := dockerCli.FilteredListNetworks(docker.NetworkFilterOpts{
		"scope": {"swarm": true},
		"label": {"com.docker.stack.namespace=" + namespace: true},
	})
	if len(networks) > 0 {
		network = &networks[0]
	}
	return
}

func createNetwork(namespace string) (network *docker.Network, err error) {
	network, err = findNetwork(namespace)
	if network != nil {
		return
	}
	network, err = dockerCli.CreateNetwork(docker.CreateNetworkOptions{
		Name:   strings.Join([]string{namespace, "Network"}, "-"),
		Driver: "overlay",
		Labels: map[string]string{
			"com.docker.stack.namespace": namespace,
		},
	})
	return
}
