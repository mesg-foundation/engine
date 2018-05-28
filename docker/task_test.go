package docker

import (
	"fmt"
	"testing"
	"time"

	"github.com/stvp/assert"
)

func TestListTasks(t *testing.T) {
	namespace := []string{"TestListTasks"}
	startTestService(namespace)
	defer StopService(namespace)
	tasks, err := ListTasks(namespace)
	assert.Nil(t, err)
	assert.NotNil(t, tasks)
	assert.Equal(t, 1, len(tasks))
	fmt.Println("number of tasks", len(tasks))
	for _, task := range tasks {
		fmt.Println("task", task.ID, task.Status, task.ServiceID)
	}
}

func TestTasksError(t *testing.T) {
	namespace := []string{"TestTasksError"}
	StartService(&ServiceOptions{
		Namespace: namespace,
		Image:     "fiifioewifewiewfifewijopwjeokpfeo",
	})
	defer StopService(namespace)

	wait := make(chan error, 1)
	errorsChan := make(chan []string, 1)
	go func() {
		for {
			errors, err := TasksError(namespace)
			if err != nil {
				wait <- err
				return
			}
			if len(errors) >= 1 {
				errorsChan <- errors
				close(wait)
				return
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	err := <-wait
	errors := <-errorsChan
	assert.Nil(t, err)
	assert.NotNil(t, errors)
	assert.Equal(t, 1, len(errors))
	assert.Equal(t, "No such image: fiifioewifewiewfifewijopwjeokpfeo:latest", errors[0])
}
