package container

import (
	"context"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
)

// CreateNetwork creates a Docker Network with a namespace
func CreateNetwork(namespace []string) (string, error) {
	network, err := FindNetwork(namespace)
	if err != nil && !docker.IsErrNotFound(err) {
		return "", err
	}
	if network.ID != "" {
		return network.ID, nil
	}
	namespaceFlat := Namespace(namespace)
	client, err := Client()
	if err != nil {
		return "", err
	}
	response, err := client.NetworkCreate(context.Background(), namespaceFlat, types.NetworkCreate{
		CheckDuplicate: true, // Cannot have 2 network with the same name
		Driver:         "overlay",
		Labels: map[string]string{
			"com.docker.stack.namespace": namespaceFlat,
		},
	})
	if err != nil {
		return "", err
	}
	return response.ID, nil
}

// DeleteNetwork deletes a Docker Network associated with a namespace
func DeleteNetwork(namespace []string) error {
	network, err := FindNetwork(namespace)
	if err != nil && !docker.IsErrNotFound(err) {
		return err
	}
	client, err := Client()
	if err != nil {
		return err
	}
	return client.NetworkRemove(context.Background(), network.ID)
}

// FindNetwork finds a Docker Network by a namespace. If no network if found, an error is returned.
func FindNetwork(namespace []string) (types.NetworkResource, error) {
	client, err := Client()
	if err != nil {
		return types.NetworkResource{}, err
	}
	return client.NetworkInspect(context.Background(), Namespace(namespace), types.NetworkInspectOptions{})
}
