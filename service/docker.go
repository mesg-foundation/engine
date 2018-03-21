package service

import (
	"os"

	docker "github.com/fsouza/go-dockerclient"
)

// DockerEndpoint is the endpoint to reach docker socket
var DockerEndpoint = "unix:///var/run/docker.sock"

var dockerCli *docker.Client

func createDockerCli(endpoint string) (client *docker.Client, err error) {
	return docker.NewClient(endpoint)
}

func init() {
	var err error
	if os.Getenv("DOCKER_HOST") != "" {
		DockerEndpoint = os.Getenv("DOCKER_HOST")
	}
	dockerCli, err = createDockerCli(DockerEndpoint)
	if err != nil {
		panic(err)
	}
}
