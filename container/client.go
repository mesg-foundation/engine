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

// Client create a docker client ready to use
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
	clientInstance = client
	return
}

func resetClient() {
	mu.Lock()
	defer mu.Unlock()
	clientInstance = nil
}

func createClient() (client *docker.Client, err error) {
	client, err = docker.NewEnvClient()
	if err != nil {
		return
	}
	client.NegotiateAPIVersion(context.Background())
	info, err := client.Info(context.Background())
	if err != nil {
		return
	}
	if info.Swarm.NodeID == "" {
		_, err = createSwarm(client)
		if err != nil {
			return
		}
	}
	err = createSharedNetworkIfNeeded(client)
	if err != nil {
		return
	}
	return
}

func createSwarm(client *docker.Client) (ID string, err error) {
	ID, err = client.SwarmInit(context.Background(), swarm.InitRequest{
		ListenAddr: "0.0.0.0:2377", // https://docs.docker.com/engine/reference/commandline/swarm_init/#usage
	})
	return
}
