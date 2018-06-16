package container

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
)

// ListTasks returns all docker tasks
func ListTasks(namespace []string) (tasks []swarm.Task, err error) {
	client, err := Client()
	if err != nil {
		return
	}
	tasks, err = client.TaskList(context.Background(), types.TaskListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "com.docker.stack.namespace=" + Namespace(namespace),
		}),
	})
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
