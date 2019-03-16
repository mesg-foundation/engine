package container

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
)

// CreateNetwork creates a Docker Network with a namespace.
func (c *DockerContainer) CreateNetwork(namespace []string) (id string, err error) {
	network, err := c.FindNetwork(namespace)
	if err != nil && !docker.IsErrNotFound(err) {
		return "", err
	}
	if network.ID != "" {
		return network.ID, nil
	}
	namespaceFlat := c.Namespace(namespace)
	response, err := c.client.NetworkCreate(context.Background(), namespaceFlat, types.NetworkCreate{
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

// DeleteNetwork deletes a Docker Network associated with a namespace.
func (c *DockerContainer) DeleteNetwork(namespace []string) error {
	network, err := c.FindNetwork(namespace)
	if docker.IsErrNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	c.client.NetworkRemove(context.Background(), network.ID)
	time.Sleep(1 * time.Second)
	return c.DeleteNetwork(namespace)
}

// FindNetwork finds a Docker Network by a namespace. If no network is found, an error is returned.
func (c *DockerContainer) FindNetwork(namespace []string) (types.NetworkResource, error) {
	return c.client.NetworkInspect(context.Background(), c.Namespace(namespace), types.NetworkInspectOptions{})
}
