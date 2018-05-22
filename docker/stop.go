package docker

import (
	"context"

	godocker "github.com/fsouza/go-dockerclient"
)

// Stop a docker service
func Stop(namespace string, dependencyName string) (err error) {
	ctx := context.Background()
	client, err := Client()
	if err != nil {
		return
	}
	if !IsRunning(namespace, dependencyName) {
		return
	}
	dockerService, err := Service(namespace, dependencyName)
	if err == nil && dockerService.ID != "" {
		err = client.RemoveService(godocker.RemoveServiceOptions{
			ID:      dockerService.ID,
			Context: ctx,
		})
	}
	return
}
