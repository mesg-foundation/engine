package container

import (
	"context"
	"errors"

	godocker "github.com/fsouza/go-dockerclient"
)

// FindContainerStrict returns a docker container if exist. Return error if not found.
func FindContainerStrict(namespace []string) (container *godocker.APIContainers, err error) {
	container, err = FindContainer(namespace)
	if err != nil {
		return
	}
	if container == nil {
		err = errors.New("Container " + Namespace(namespace) + " not found")
	}
	return
}

// FindContainer returns a docker container if exist
func FindContainer(namespace []string) (container *godocker.APIContainers, err error) {
	client, err := Client()
	if err != nil {
		return
	}
	containers, err := client.ListContainers(godocker.ListContainersOptions{
		Context: context.Background(),
		Limit:   1,
		Filters: map[string][]string{
			"label": []string{
				"com.docker.stack.namespace=" + Namespace(namespace),
			},
		},
	})
	if err != nil {
		return
	}
	if len(containers) == 1 {
		container = &containers[0]
	}
	return
}

// ContainerStatus returns the status of a docker container
func ContainerStatus(namespace []string) (status StatusType, err error) {
	container, err := FindContainer(namespace)
	if err != nil {
		return
	}
	status = STOPPED
	if container != nil && container.State == "running" {
		status = RUNNING
	}
	return
}
