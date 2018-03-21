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
	info, err := client.Info()
	if info.Swarm.NodeID != "" {
		return
	}
	// TODO remove and make CI works
	if os.Getenv("CI") == "true" {
		return
	}
	res, err := client.InitSwarm(docker.InitSwarmOptions{
		Context: context.Background(),
		InitRequest: swarm.InitRequest{
			ListenAddr: "0.0.0.0:2377", // https://docs.docker.com/engine/reference/commandline/swarm_init/#usage
		},
	})
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
