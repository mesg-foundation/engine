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

// TasksError returns the error of matching tasks
func TasksError(namespace []string) (errors []string, err error) {
	tasks, err := ListTasks(namespace)
	if err != nil {
		return
	}
	for _, task := range tasks {
		if task.Status.Err != "" {
			errors = append(errors, task.Status.Err)
		}
	}
	return
}
