package container

import (
	"errors"
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/stvp/assert"
)

func TestListTasks(t *testing.T) {
	namespace := []string{"namespace"}
	tasks := []swarm.Task{
		{ID: "1"},
		{ID: "2"},
	}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideTaskList(tasks, nil)

	tasks1, err := c.ListTasks(namespace)
	assert.Nil(t, err)
	assert.Equal(t, tasks, tasks1)
	assert.Equal(t, len(tasks), len(tasks1))
}

var errTaskList = errors.New("task list")

func TestListTasksError(t *testing.T) {
	namespace := []string{"namespace"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideTaskList(nil, errTaskList)

	_, err := c.ListTasks(namespace)
	assert.Equal(t, errTaskList, err)
}

func TestTasksError(t *testing.T) {
	namespace := []string{"namespace"}
	tasks := []swarm.Task{
		{
			ID:     "1",
			Status: swarm.TaskStatus{Err: "1-err"},
		},
		{
			ID:     "1",
			Status: swarm.TaskStatus{Err: "2-err"},
		},
	}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideTaskList(tasks, nil)

	errors, err := c.TasksError(namespace)
	assert.Nil(t, err)
	assert.Equal(t, len(tasks), len(errors))
	assert.Equal(t, tasks[0].Status.Err, errors[0])
	assert.Equal(t, tasks[1].Status.Err, errors[1])
}
