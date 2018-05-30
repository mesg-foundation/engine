package container

import (
	"context"
	"sync"

	"github.com/docker/docker/api/types/swarm"
	docker "github.com/fsouza/go-dockerclient"
)

var sharedNetworkNamespace = []string{"shared"}

var dockerCliInstance *docker.Client
var mu sync.Mutex

// Client create a docker client ready to use
func Client() (client *docker.Client, err error) {
	mu.Lock()
	defer mu.Unlock()
	if dockerCliInstance != nil {
		client = dockerCliInstance
		return
	}
	client, err = createClient()
	if err != nil {
		return
	}
	dockerCliInstance = client
	return
}

func resetClient() {
	mu.Lock()
	defer mu.Unlock()
	dockerCliInstance = nil
}

func createClient() (client *docker.Client, err error) {
	client, err = docker.NewClientFromEnv()
	if err != nil {
		return
	}
	info, err := client.Info()
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
	ID, err = client.InitSwarm(docker.InitSwarmOptions{
		Context: context.Background(),
		InitRequest: swarm.InitRequest{
			ListenAddr: "0.0.0.0:2377", // https://docs.docker.com/engine/reference/commandline/swarm_init/#usage
		},
	})
	return
}
