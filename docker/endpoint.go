package docker

import (
	"strings"

	godocker "github.com/fsouza/go-dockerclient"
)

// FindEndpoint returns the endoint of a docker container in a specific network
func FindEndpoint(networkName string, containerName string) (endpoint *godocker.Endpoint, err error) {
	network, err := FindNetwork(networkName)
	if network == nil || err != nil {
		return
	}
	namespace := Namespace([]string{containerName})
	for _, e := range network.Containers {
		if strings.Contains(e.Name, namespace) {
			endpoint = &e
			break
		}
	}
	return
}

// FindIP returns the ipv4 of a docker container in a specific network
func FindIP(networkName string, containerName string) (IP string, err error) {
	endpoint, err := FindEndpoint(networkName, containerName)
	if endpoint == nil || err != nil {
		return
	}
	split := strings.Split(endpoint.IPv4Address, "/")
	if len(split) != 2 {
		return
	}
	IP = split[0]
	return
}
