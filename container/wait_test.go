package container

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/mesg-foundation/core/utils/docker/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TODO: support all status types.
func mockWaitForStatus(t *testing.T, m *mocks.CommonAPIClient, namespace string, wantedStatus StatusType) {
	var (
		containerID = "1"
	)

	m.On("TaskList", mock.Anything, types.TaskListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "com.docker.stack.namespace=" + namespace,
		}),
	}).Once().
		Return([]swarm.Task{}, nil)
	m.On("ContainerList", mock.AnythingOfType("*context.timerCtx"), types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "com.docker.stack.namespace=" + namespace,
		}),
		Limit: 1,
	}).Once().
		Return([]types.Container{{ID: containerID}}, nil)

	containerInspect := m.On("ContainerInspect", mock.AnythingOfType("*context.timerCtx"), containerID).Once()
	serviceInspect := m.On("ServiceInspectWithRaw", mock.Anything, namespace, types.ServiceInspectOptions{}).Once()
	switch wantedStatus {
	case STOPPED:
		containerInspect.Return(types.ContainerJSON{}, dockertest.NotFoundErr{})
		serviceInspect.Return(swarm.Service{}, nil, dockertest.NotFoundErr{})
	default:
		t.Errorf("unhandled status %v", wantedStatus)
	}
}

func TestWaitForStatusRunning(t *testing.T) {
	namespace := []string{"namespace"}
	containerID := "1"
	containerData := []types.Container{
		{ID: containerID},
	}
	containerJSONData := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID:    containerID,
			State: &types.ContainerState{Running: true},
		},
	}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideContainerList(containerData, nil)
	dt.ProvideContainerInspect(containerJSONData, nil)

	require.Nil(t, c.waitForStatus(namespace, RUNNING))
}

func TestWaitForStatusStopped(t *testing.T) {
	namespace := []string{"namespace"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideServiceInspectWithRaw(swarm.Service{}, nil, dockertest.NotFoundErr{})
	dt.ProvideContainerInspect(types.ContainerJSON{}, dockertest.NotFoundErr{})
	require.Nil(t, c.waitForStatus(namespace, STOPPED))
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

	require.Equal(t, "1-err, 2-err", c.waitForStatus(namespace, RUNNING).Error())

	select {
	case <-dt.LastContainerList():
		t.Error("container list shouldn't be called")
	default:
	}
}
