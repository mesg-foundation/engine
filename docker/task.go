package docker

import (
	"context"

	"github.com/docker/docker/api/types/swarm"
	godocker "github.com/fsouza/go-dockerclient"
)

// ListTasks returns all docker tasks
func ListTasks(namespace []string) (tasks []swarm.Task, err error) {
	client, err := Client()
	if err != nil {
		return
	}
	tasks, err = client.ListTasks(godocker.ListTasksOptions{
		Context: context.Background(),
		Filters: map[string][]string{
			"service": []string{Namespace(namespace)},
		},
	})
	if err != nil {
		return
	}
	return
}
