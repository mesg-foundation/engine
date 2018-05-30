package container

import (
	"context"

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

// ContainerStatus returns the status of a docker container
func ContainerStatus(namespace []string) (status StatusType, err error) {
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

// func WaitContainer(namespace []string, waitCondition containerTypes.WaitCondition, timeout time.Duration) {
// 	client, err := Client()
// 	if err != nil {
// 		return
// 	}
// 	container, err := FindContainer(namespace)
// 	if err != nil {
// 		return
// 	}
// 	ctx, cancel := context.WithTimeout(context.Background(), timeout)
// 	defer cancel()
// 	ok, errC := client.ContainerWait(ctx, container.ID, waitCondition)
// 	select {
// 	case <-ok:
// 		return
// 	case err = <-errC:
// 		return
// 	}
// }
