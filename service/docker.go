package service

import (
	"context"
	"errors"
	"sync"

	"github.com/docker/docker/api/types/swarm"
	docker "github.com/fsouza/go-dockerclient"
)

var dockerCliInstance *docker.Client
var mu sync.Mutex

// DockerClient returns a preconfigured docker client
func DockerClient() (client *docker.Client, err error) {
	return dockerCli()
}

func dockerCli() (client *docker.Client, err error) {
	mu.Lock()
	defer mu.Unlock()
	if dockerCliInstance != nil {
		client = dockerCliInstance
		return
	}
	client, err = createDockerCli()
	if err != nil {
		return nil, err
	}
	dockerCliInstance = client
	return
}

func resetCliInstance() {
	mu.Lock()
	defer mu.Unlock()
	dockerCliInstance = nil
}

func createDockerCli() (client *docker.Client, err error) {
	client, err = docker.NewClientFromEnv()
	if err != nil {
		return
	}

	info, errSwarm := client.Info()
	if errSwarm != nil || info.Swarm.NodeID != "" {
		err = createNetworkIfNeeded(client)
		return
	}
	_, err = createSwarm(client)
	if err != nil {
		return
	}

	err = createNetworkIfNeeded(client)

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

func createNetworkIfNeeded(client *docker.Client) (err error) {
	_, findError := SharedNetwork(client)
	if findError != nil {
		// Create the new network needed to run containers
		_, err = client.CreateNetwork(docker.CreateNetworkOptions{
			Context:        context.Background(),
			CheckDuplicate: true,
			Name:           "daemon-network",
			Driver:         "overlay",
		})
	}
	return
}

// SharedNetwork returns the shared network created to connect services and daemon
func SharedNetwork(client *docker.Client) (network docker.Network, err error) {
	networks, err := client.FilteredListNetworks(docker.NetworkFilterOpts{
		"name": {"daemon-network": true},
	})
	if err != nil {
		return
	}
	if len(networks) == 0 {
		err = errors.New("Cannot find the appropriate docker network")
		return
	}
	network = networks[0]
	return
}
