package docker

import (
	"context"
	"fmt"
	"sync"

	"github.com/docker/docker/api/types/swarm"
	godocker "github.com/fsouza/go-dockerclient"
	"github.com/logrusorgru/aurora"
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
		return nil, err
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
	ID, err := createSwarm(client)
	if err == nil {
		fmt.Println(aurora.Green("Docker swarm node created"), ID)
	}
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
