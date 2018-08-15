// +build integration

package container

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIntegrationListTasks(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
	namespace := []string{"TestListTasks"}
	startTestService(namespace)
	defer c.StopService(namespace)
	tasks, err := c.ListTasks(namespace)
	require.Nil(t, err)
	require.NotNil(t, tasks)
	require.Equal(t, 1, len(tasks))
}

func TestIntegrationTasksError(t *testing.T) {
	c, err := New()
	require.Nil(t, err)
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
	require.Nil(t, err)
	require.NotNil(t, errors)
	require.True(t, len(errors) > 0)
	require.Equal(t, "No such image: fiifioewifewiewfifewijopwjeokpfeo:latest", errors[0])
}
