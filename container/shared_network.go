package container

import (
	"context"

	"github.com/docker/docker/api/types"
	docker "github.com/docker/docker/client"
)

func createSharedNetworkIfNeeded(client *docker.Client) (err error) {
	network, err := sharedNetwork(client)
	if docker.IsErrNotFound(err) {
		err = nil
	}
	if err != nil {
		return
	}
	if network.ID != "" {
		return
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
	return
}

// sharedNetwork returns the shared network created to connect services and daemon
func sharedNetwork(client *docker.Client) (network types.NetworkResource, err error) {
	network, err = client.NetworkInspect(context.Background(), Namespace(sharedNetworkNamespace), types.NetworkInspectOptions{})
	return
}

// SharedNetworkID returns the id of the shared network
func SharedNetworkID() (networkID string, err error) {
	client, err := Client()
	if err != nil {
		return
	}
	network, err := sharedNetwork(client)
	if err != nil {
		return
	}
	networkID = network.ID
	return
}
