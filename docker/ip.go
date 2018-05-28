package docker

import (
	"strings"
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
func FindServiceIP(networkNamespace []string, serviceNamespace []string) (IP string, err error) {
	network, err := FindNetwork(networkNamespace)
	if network == nil || err != nil {
		return
	}
	service, err := FindService(serviceNamespace)
	if service == nil || err != nil {
		return
	}
	for _, virtualIP := range service.Endpoint.VirtualIPs {
		if virtualIP.NetworkID == network.ID {
			IP = removeIPSuffix(virtualIP.Addr)
			return
		}
	}
	return
}

// FindContainerIP returns the ipv4 of a docker container in a specific network
func FindContainerIP(networkNamespace []string, containerNamespace []string) (IP string, err error) {
	network, err := FindNetwork(networkNamespace)
	if network == nil || err != nil {
		return
	}
	for _, endpoint := range network.Containers {
		if strings.Contains(endpoint.Name, Namespace(containerNamespace)) {
			IP = removeIPSuffix(endpoint.IPv4Address)
			break
		}
	}
	return
}
