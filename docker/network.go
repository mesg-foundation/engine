package docker

import (
	godocker "github.com/fsouza/go-dockerclient"
)

const networkNamespacePrefix string = "network"

func networkNamespace(name string) string {
	return Namespace([]string{networkNamespacePrefix, name})
}

// CreateNetwork creates a Docker Network with a namespace
func CreateNetwork(name string) (network *godocker.Network, err error) {
	network, err = FindNetwork(name)
	if network != nil || err != nil {
		return
	}
	namespace := networkNamespace(name)
	client, err := Client()
	if err != nil {
		return
	}
	network, err = client.CreateNetwork(godocker.CreateNetworkOptions{
		Name:           namespace,
		CheckDuplicate: true, // Cannot have 2 network with the same name
		Driver:         "overlay",
		Labels: map[string]string{
			"com.docker.stack.namespace": namespace,
		},
	})
	return
}

// DeleteNetwork deletes a Docker Network associated with a namespace
func DeleteNetwork(name string) (err error) {
	network, err := FindNetwork(name)
	if network == nil || err != nil {
		return
	}
	client, err := Client()
	if err != nil {
		return
	}
	return client.RemoveNetwork(network.ID)
}

// FindNetwork finds a Docker Network by a namespace
func FindNetwork(name string) (network *godocker.Network, err error) {
	client, err := Client()
	if err != nil {
		return
	}
	namespace := networkNamespace(name)
	return client.NetworkInfo(namespace)
}
