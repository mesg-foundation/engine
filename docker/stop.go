package docker

import (
	"context"

	godocker "github.com/fsouza/go-dockerclient"
)

// Stop a docker service
func Stop(name []string) (err error) {
	client, err := Client()
	if err != nil {
		return
	}
	if !IsRunning(name) {
		return
	}
	dockerService, err := FindService(name)
	if err == nil && dockerService.ID != "" {
		err = client.RemoveService(godocker.RemoveServiceOptions{
			ID:      dockerService.ID,
			Context: context.Background(),
		})
	}
	return
}
