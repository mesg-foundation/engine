package service

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/docker/docker/api/types/swarm"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/logrusorgru/aurora"
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
	var wg sync.WaitGroup
	wg.Add(2)

	// Swarm check
	go func() {
		defer wg.Done()
		info, err := client.Info()
		if err != nil || info.Swarm.NodeID != "" {
			return
		}
		ID, err := createSwarm(client)
		if err == nil {
			fmt.Println(aurora.Green("Docker swarm node created"), ID)
		} else {
			return
		}
	}()

	// Network check
	go func() {
		defer wg.Done()
		network, e := SharedNetwork(client)
		if e == nil && network.ID != "" {
			return
		}
		// Create the new network needed to run containers
		_, err = client.CreateNetwork(docker.CreateNetworkOptions{
			Context:        context.Background(),
			CheckDuplicate: true,
			Name:           "daemon-network",
			Driver:         "overlay",
		})
	}()

	wg.Wait()
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
