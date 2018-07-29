package container

import (
	"testing"
	"time"

	"github.com/stvp/assert"
)

func TestListTasks(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestListTasks"}
	startTestService(namespace)
	defer c.StopService(namespace)
	tasks, err := c.ListTasks(namespace)
	assert.Nil(t, err)
	assert.NotNil(t, tasks)
	assert.Equal(t, 1, len(tasks))
}

func TestTasksError(t *testing.T) {
	c, err := New()
	assert.Nil(t, err)
	namespace := []string{"TestTasksError"}
	c.StartService(ServiceOptions{
		Image:     "fiifioewifewiewfifewijopwjeokpfeo",
		Namespace: namespace,
	})
	defer c.StopService(namespace)
	var errors []string
	for {
		errors, err = c.TasksError(namespace)
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
