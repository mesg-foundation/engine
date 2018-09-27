package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
)

// CreateNetwork creates a Docker Network with a namespace.
func (c *Container) CreateNetwork(namespace []string) (id string, err error) {
	network, err := c.FindNetwork(namespace)
	if err != nil && !docker.IsErrNotFound(err) {
		return "", err
	}
	if network.ID != "" {
		return network.ID, nil
	}
	namespaceFlat := c.Namespace(namespace)
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	response, err := c.client.NetworkCreate(ctx, namespaceFlat, types.NetworkCreate{
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
func (c *Container) DeleteNetwork(namespace []string) error {
	network, err := c.FindNetwork(namespace)
	if docker.IsErrNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	messageChan, errChan := c.client.Events(ctx, types.EventsOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "network",
			Value: network.ID,
		}, filters.KeyValuePair{
			Key:   "event",
			Value: "destroy",
		}),
	})
	err = c.client.NetworkRemove(ctx, network.ID)
	if err != nil {
		return err
	}
	select {
	case <-messageChan:
		return nil
	case err := <-errChan:
		return err
	}
}

// FindNetwork finds a Docker Network by a namespace. If no network is found, an error is returned.
func (c *Container) FindNetwork(namespace []string) (types.NetworkResource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	return c.client.NetworkInspect(ctx, c.Namespace(namespace), types.NetworkInspectOptions{})
}
