package docker

import (
	"context"
	"sync"

	"github.com/docker/docker/api/types/swarm"
	godocker "github.com/fsouza/go-dockerclient"
)

var dockerCliInstance *godocker.Client
var mu sync.Mutex

// Client create a docker client ready to use
func Client() (client *godocker.Client, err error) {
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

func createClient() (client *godocker.Client, err error) {
	client, err = godocker.NewClientFromEnv()
	if err != nil {
		return
	}
	info, err := client.Info()
	if err != nil || info.Swarm.NodeID != "" {
		return
	}
	_, err = createSwarm(client)
	return
}

func createSwarm(client *godocker.Client) (ID string, err error) {
	ID, err = client.InitSwarm(godocker.InitSwarmOptions{
		Context: context.Background(),
		InitRequest: swarm.InitRequest{
			ListenAddr: "0.0.0.0:2377", // https://docs.docker.com/engine/reference/commandline/swarm_init/#usage
		},
	})
	return
}
