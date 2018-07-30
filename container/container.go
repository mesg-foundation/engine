package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
)

// FindContainer returns a docker container if exist
func FindContainer(namespace []string) (types.ContainerJSON, error) {
	client, err := Client()
	if err != nil {
		return types.ContainerJSON{}, err
	}
	containers, err := client.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "com.docker.stack.namespace=" + Namespace(namespace),
		}),
		Limit: 1,
	})
	if err != nil {
		return types.ContainerJSON{}, err
	}
	containerID := ""
	if len(containers) == 1 {
		containerID = containers[0].ID
	}
	return client.ContainerInspect(context.Background(), containerID)
}

// Status returns the status of a docker container
func Status(namespace []string) (StatusType, error) {
	status := STOPPED
	container, err := FindContainer(namespace)
	if docker.IsErrNotFound(err) {
		return status, nil
	}
	if err != nil {
		return status, err
	}
	if container.State.Running {
		status = RUNNING
	}
	return status, nil
}
