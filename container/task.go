package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
)

// ListTasks returns all docker tasks
func ListTasks(namespace []string) ([]swarm.Task, error) {
	client, err := Client()
	if err != nil {
		return nil, err
	}
	return client.TaskList(context.Background(), types.TaskListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "com.docker.stack.namespace=" + Namespace(namespace),
		}),
	})
}

// TasksError returns the error of matching tasks
func TasksError(namespace []string) ([]string, error) {
	tasks, err := ListTasks(namespace)
	if err != nil {
		return nil, err
	}
	var errors []string
	for _, task := range tasks {
		if task.Status.Err != "" {
			errors = append(errors, task.Status.Err)
		}
	}
	return errors, nil
}
