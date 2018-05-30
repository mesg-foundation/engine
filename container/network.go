package container

import (
	godocker "github.com/fsouza/go-dockerclient"
)

// CreateNetwork creates a Docker Network with a namespace
func CreateNetwork(namespace []string) (network *godocker.Network, err error) {
	network, err = FindNetwork(namespace)
	if network != nil || err != nil {
		return
	}
	namespaceFlat := Namespace(namespace)
	client, err := Client()
	if err != nil {
		return
	}
	network, err = client.CreateNetwork(godocker.CreateNetworkOptions{
		Name:           namespaceFlat,
		CheckDuplicate: true, // Cannot have 2 network with the same name
		Driver:         "overlay",
		Labels: map[string]string{
			"com.docker.stack.namespace": namespaceFlat,
		},
	})
	return
}

// DeleteNetwork deletes a Docker Network associated with a namespace
func DeleteNetwork(namespace []string) (err error) {
	network, err := FindNetwork(namespace)
	if network == nil || err != nil {
		return
	}
	client, err := Client()
	if err != nil {
		return
	}
	return client.RemoveNetwork(network.ID)
}

// FindNetworkStrict finds a Docker Network by a namespace. If no network if found, an error is returned.
func FindNetworkStrict(namespace []string) (network *godocker.Network, err error) {
	client, err := Client()
	if err != nil {
		return
	}
	return client.NetworkInfo(Namespace(namespace))
}

// FindNetwork finds a Docker Network by a namespace. If no network if found, NO error is returned.
func FindNetwork(namespace []string) (network *godocker.Network, err error) {
	network, err = FindNetworkStrict(namespace)
	if err != nil {
		switch err.(type) {
		case *godocker.NoSuchNetwork:
			err = nil
		default:
		}
	}
	return
}
