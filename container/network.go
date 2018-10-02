package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
)

// EventType is a type to define the kind of event related to the network
type EventType string

// List of Network status to listen to confirm a network deletion
const (
	EventRemove  EventType = "remove"
	EventDestroy           = "destroy"
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
// event parameter can be "destroy" or "remove". If the network was used by a service, the event to use is "destroy". If the network has not been used, the event is "remove".
// Remove removes the reference from Docker to the network.
// Destroy removes the network from Docker active network.
func (c *DockerContainer) DeleteNetwork(namespace []string, event EventType) error {
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
			Value: string(event),
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
func (c *DockerContainer) FindNetwork(namespace []string) (types.NetworkResource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	return c.client.NetworkInspect(ctx, c.Namespace(namespace), types.NetworkInspectOptions{})
}
