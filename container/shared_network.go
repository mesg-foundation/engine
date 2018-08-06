package container

import (
	"context"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
)

func createSharedNetworkIfNeeded(client *docker.Client) error {
	network, err := sharedNetwork(client)
	if err != nil && !docker.IsErrNotFound(err) {
		return err
	}
	if network.ID != "" {
		return nil
	}
	// Create the new network needed to run containers
	namespace := Namespace(sharedNetworkNamespace)
	_, err = client.NetworkCreate(context.Background(), namespace, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "overlay",
		Labels: map[string]string{
			"com.docker.stack.namespace": namespace,
		},
	})
	return err
}

// sharedNetwork returns the shared network created to connect services and MESG Core.
func sharedNetwork(client *docker.Client) (types.NetworkResource, error) {
	return client.NetworkInspect(context.Background(), Namespace(sharedNetworkNamespace), types.NetworkInspectOptions{})
}

// SharedNetworkID returns the ID of the shared network.
func SharedNetworkID() (string, error) {
	client, err := Client()
	if err != nil {
		return "", err
	}
	network, err := sharedNetwork(client)
	if err != nil {
		return "", err
	}
	return network.ID, nil
}
