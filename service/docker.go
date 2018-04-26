package service

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/swarm"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/logrusorgru/aurora"
)

var dockerCliInstance *docker.Client

func dockerCli() (client *docker.Client, err error) {
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

func createDockerCli() (client *docker.Client, err error) {
	client, err = docker.NewClientFromEnv()
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

func createSwarm(client *docker.Client) (ID string, err error) {
	ID, err = client.InitSwarm(docker.InitSwarmOptions{
		Context: context.Background(),
		InitRequest: swarm.InitRequest{
			ListenAddr: "0.0.0.0:2377", // https://docs.docker.com/engine/reference/commandline/swarm_init/#usage
		},
	})
	return
}
