package container

import (
	"context"
	"sync"

	"github.com/docker/docker/api/types/swarm"
	docker "github.com/docker/docker/client"
)

var sharedNetworkNamespace = []string{"shared"}

var clientInstance *docker.Client
var mu sync.Mutex

// Client creates a ready to use docker client.
func Client() (client *docker.Client, err error) {
	mu.Lock()
	defer mu.Unlock()
	if clientInstance != nil {
		client = clientInstance
		return
	}
	client, err = createClient()
	if err != nil {
		return
	}
	err = createSwarmIfNeeded(client)
	if err != nil {
		return
	}
	err = createSharedNetworkIfNeeded(client)
	if err != nil {
		return
	}
	clientInstance = client
	return
}

func resetClient() {
	mu.Lock()
	defer mu.Unlock()
	clientInstance = nil
}

func createClient() (*docker.Client, error) {
	client, err := docker.NewEnvClient()
	if err != nil {
		return nil, err
	}
	client.NegotiateAPIVersion(context.Background())
	return client, nil
}

func createSwarmIfNeeded(client *docker.Client) error {
	info, err := client.Info(context.Background())
	if err != nil {
		return err
	}
	if info.Swarm.NodeID != "" {
		return nil
	}
	_, err = client.SwarmInit(context.Background(), swarm.InitRequest{
		ListenAddr: "0.0.0.0:2377", // https://docs.docker.com/engine/reference/commandline/swarm_init/#usage
	})
	return err
}
