package container

import (
	"context"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
)

// CreateNetwork creates a Docker Network with a namespace
func CreateNetwork(namespace []string) (networkID string, err error) {
	network, err := FindNetwork(namespace)
	if docker.IsErrNotFound(err) {
		err = nil
	}
	if err != nil {
		return
	}
	if network.ID != "" {
		networkID = network.ID
		return
	}
	namespaceFlat := Namespace(namespace)
	client, err := Client()
	if err != nil {
		return
	}
	response, err := client.NetworkCreate(context.Background(), namespaceFlat, types.NetworkCreate{
		CheckDuplicate: true, // Cannot have 2 network with the same name
		Driver:         "overlay",
		Labels: map[string]string{
			"com.docker.stack.namespace": namespaceFlat,
		},
	})
	if err != nil {
		return
	}
	networkID = response.ID
	return
}

// DeleteNetwork deletes a Docker Network associated with a namespace
func DeleteNetwork(namespace []string) (err error) {
	network, err := FindNetwork(namespace)
	if docker.IsErrNotFound(err) {
		err = nil
		return
	}
	if err != nil {
		return
	}
	client, err := Client()
	if err != nil {
		return
	}
	err = client.NetworkRemove(context.Background(), network.ID)
	return
}

// FindNetwork finds a Docker Network by a namespace. If no network if found, an error is returned.
func FindNetwork(namespace []string) (network types.NetworkResource, err error) {
	client, err := Client()
	if err != nil {
		return
	}
	network, err = client.NetworkInspect(context.Background(), Namespace(namespace), types.NetworkInspectOptions{})
	return
}
