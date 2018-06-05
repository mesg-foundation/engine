package container

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	docker "github.com/docker/docker/client"
)

// FindContainer returns a docker container if exist
func FindContainer(namespace []string) (container types.ContainerJSON, err error) {
	client, err := Client()
	if err != nil {
		return
	}
	containers, err := client.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "com.docker.stack.namespace=" + Namespace(namespace),
		}),
		Limit: 1,
	})
	if err != nil {
		return
	}
	containerID := ""
	if len(containers) == 1 {
		containerID = containers[0].ID
	}
	container, err = client.ContainerInspect(context.Background(), containerID)
	return
}

// Status returns the status of a docker container
func Status(namespace []string) (status StatusType, err error) {
	status = STOPPED
	container, err := FindContainer(namespace)
	if docker.IsErrNotFound(err) {
		err = nil
		return
	}
	if err != nil {
		return
	}
	if container.State.Running {
		status = RUNNING
	}
	return
}

// WaitForContainerStatus wait for the container to have the provided status until it reach the timeout
func WaitForContainerStatus(namespace []string, status StatusType) (err error) {
	for {
		var currentStatus StatusType
		currentStatus, err = Status(namespace)
		if err != nil {
			break
		}
		if currentStatus == status {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	return
}
