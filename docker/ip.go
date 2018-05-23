package docker

import (
	"strings"

	godocker "github.com/fsouza/go-dockerclient"
)

// removeIPSuffix removes for eg the "/24" of "10.0.0.45/24". Return the ip if no suffix.
func removeIPSuffix(IP string) string {
	split := strings.Split(IP, "/")
	if len(split) != 2 {
		return IP
	}
	return split[0]
}

// FindServiceIP returns the IP of a docker service in a specific network
func FindServiceIP(networkName string, serviceName string) (IP string, err error) {
	network, err := FindNetwork(networkName)
	if network == nil || err != nil {
		return
	}
	service, err := FindService([]string{serviceName})
	if err != nil {
		return
	}
	for _, virtualIP := range service.Endpoint.VirtualIPs {
		if virtualIP.NetworkID == network.ID {
			IP = removeIPSuffix(virtualIP.Addr)
			break
		}
	}
	return
}

// FindContainerIP returns the ipv4 of a docker container in a specific network
func FindContainerIP(networkName string, containerName string) (IP string, err error) {
	endpoint, err := findContainerEndpoint(networkName, containerName)
	if endpoint == nil || err != nil {
		return
	}
	IP = removeIPSuffix(endpoint.IPv4Address)
	return
}

// findContainerEndpoint returns the endpoint of a docker container in a specific network
func findContainerEndpoint(networkName string, containerName string) (endpoint *godocker.Endpoint, err error) {
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
