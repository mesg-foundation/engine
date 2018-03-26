package service

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types/swarm"
	"github.com/logrusorgru/aurora"

	docker "github.com/fsouza/go-dockerclient"
)

// DockerEndpoint is the endpoint to reach docker socket
var DockerEndpoint = "unix:///var/run/docker.sock"

var dockerCli *docker.Client

func createDockerCli(endpoint string) (client *docker.Client, err error) {
	client, err = docker.NewClient(endpoint)
	if err != nil {
		return
	}
	// TODO remove and make CI works
	if os.Getenv("CI") == "true" {
		return
	}
	info, err := client.Info()
	if err != nil {
		return
	}
	if info.Swarm.NodeID != "" {
		return
	}
	res, err := client.InitSwarm(docker.InitSwarmOptions{
		Context: context.Background(),
		InitRequest: swarm.InitRequest{
			ListenAddr: "0.0.0.0:2377", // https://docs.docker.com/engine/reference/commandline/swarm_init/#usage
		},
	})
	if err != nil {
		return
	}
	fmt.Println(aurora.Green("Docker swarm node created"), res)
	return
}

func init() {
	var err error
	dockerCli, err = createDockerCli(DockerEndpoint)
	if err != nil {
		panic(err)
	}
}
