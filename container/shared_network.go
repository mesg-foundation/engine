package container

import (
	"context"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
)

// SharedNetworkID returns the ID of the shared network.
func (c *DockerContainer) SharedNetworkID() (networkID string, err error) {
	network, err := c.sharedNetwork()
	if err != nil {
		return "", err
	}
	return network.ID, nil
}

func (c *DockerContainer) createSharedNetworkIfNeeded() error {
	network, err := c.sharedNetwork()
	if err != nil && !docker.IsErrNotFound(err) {
		return err
	}
	if network.ID != "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()

	// Create the new network needed to run containers.
	namespace := c.Namespace([]string{})
	_, err = c.client.NetworkCreate(ctx, namespace, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "overlay",
		Labels: map[string]string{
			"com.docker.stack.namespace": namespace,
		},
	})
	return err
}

// sharedNetwork returns the shared network created to connect services and MESG Core.
func (c *DockerContainer) sharedNetwork() (network types.NetworkResource, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	return c.client.NetworkInspect(ctx, c.Namespace([]string{}), types.NetworkInspectOptions{})
}
