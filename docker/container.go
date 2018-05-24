package docker

import (
	"context"
	"errors"
	"time"

	godocker "github.com/fsouza/go-dockerclient"
)

// FindContainer returns a running docker container if exist
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

// WaitForContainer wait for the container to run until it reach the timeout
func WaitForContainer(namespace []string) (err error) {
	maxWait := 20 * time.Second
	start := time.Now()
	var status StatusType
	for {
		status, err = ContainerStatus(namespace)
		if err != nil {
			return
		}
		if status == RUNNING {
			return
		}

		diff := time.Now().Sub(start)
		if diff.Nanoseconds() >= int64(maxWait) {
			err = errors.New("Wait too long for the container, timeout reached")
			return
		}

		time.Sleep(500 * time.Millisecond)
	}
}
