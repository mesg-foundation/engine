package container

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/stvp/assert"
)

func TestWaitForStatusRunning(t *testing.T) {
	namespace := []string{"namespace"}
	containerID := "1"
	containerData := types.Container{ID: containerID}
	containerJSONData := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID:    containerID,
			State: &types.ContainerState{Running: true},
		},
	}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideContainer(containerData)
	dt.ProvideContainerInspect(containerJSONData)

	assert.Nil(t, c.waitForStatus(namespace, RUNNING))
}

func TestWaitForStatusStopped(t *testing.T) {
	namespace := []string{"namespace"}
	containerID := "1"
	containerData := types.Container{ID: containerID}
	containerJSONData := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID:    containerID,
			State: &types.ContainerState{},
		},
	}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideContainer(containerData)
	dt.ProvideContainerInspect(containerJSONData)

	assert.Nil(t, c.waitForStatus(namespace, STOPPED))
}

func TestWaitForStatusTaskError(t *testing.T) {
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

	assert.Equal(t, "1-err, 2-err", c.waitForStatus(namespace, RUNNING).Error())

	select {
	case <-dt.LastContainerList():
		t.Error("container list shouldn't be called")
	default:
	}
}
