package container

import (
	docker "github.com/fsouza/go-dockerclient"
)

func createSharedNetworkIfNeeded(client *docker.Client) (err error) {
	_, err = sharedNetwork(client)
	if err != nil {
		switch err.(type) {
		case *docker.NoSuchNetwork:
			namespace := Namespace(sharedNetworkNamespace)
			// Create the new network needed to run containers
			_, err = client.CreateNetwork(docker.CreateNetworkOptions{
				Name:           namespace,
				CheckDuplicate: true,
				Driver:         "overlay",
				Labels: map[string]string{
					"com.docker.stack.namespace": namespace,
				},
			})
		default:
		}
	}
	return
}

// sharedNetwork returns the shared network created to connect services and daemon
func sharedNetwork(client *docker.Client) (network *docker.Network, err error) {
	network, err = client.NetworkInfo(Namespace(sharedNetworkNamespace))
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
