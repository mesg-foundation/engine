package container

import (
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
}

func TestTasksError(t *testing.T) {
	namespace := []string{"TestTasksError"}
	StartService(ServiceOptions{
		Image:     "fiifioewifewiewfifewijopwjeokpfeo",
		Namespace: namespace,
	})
	defer StopService(namespace)
	var errors []string
	var err error
	for {
		errors, err = TasksError(namespace)
		if err != nil || len(errors) > 0 {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	assert.Nil(t, err)
	assert.NotNil(t, errors)
	assert.True(t, len(errors) > 0)
	assert.Equal(t, "No such image: fiifioewifewiewfifewijopwjeokpfeo:latest", errors[0])
}
