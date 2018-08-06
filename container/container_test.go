package container

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"

	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/stvp/assert"
)

func TestNew(t *testing.T) {
	dt := dockertest.New()
	c, err := New(ClientOption(dt.Client()))
	assert.Nil(t, err)
	assert.NotNil(t, c)

	select {
	case <-dt.LastNegotiateAPIVersion():
	default:
		t.Fatal("should negotiate api version")
	}

	select {
	case <-dt.LastInfo():
	default:
		t.Error("should fetch info")
	}

	assert.Equal(t, "0.0.0.0:2377", (<-dt.LastSwarmInit()).Request.ListenAddr)

	ln := <-dt.LastNetworkCreate()
	assert.Equal(t, "mesg-shared", ln.Name)
	assert.Equal(t, types.NetworkCreate{
		CheckDuplicate: true,
		Driver:         "overlay",
		Labels: map[string]string{
			"com.docker.stack.namespace": ln.Name,
		},
	}, ln.Options)
}

func TestNewWithExistingNode(t *testing.T) {
	dt := dockertest.New()
	dt.ProvideInfo(types.Info{Swarm: swarm.Info{NodeID: "1"}}, nil)

	c, err := New(ClientOption(dt.Client()))
	assert.Nil(t, err)
	assert.NotNil(t, c)

	select {
	case <-dt.LastSwarmInit():
		t.Fail()
	default:
	}
}

func TestFindContainerNonExistent(t *testing.T) {
	namespace := []string{"namespace"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideContainerList(nil, dockertest.NotFoundErr{})

	_, err := c.FindContainer(namespace)
	assert.Equal(t, dockertest.NotFoundErr{}, err)

	assert.Equal(t, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "com.docker.stack.namespace=" + Namespace(namespace),
		}),
		Limit: 1,
	}, (<-dt.LastContainerList()).Options)
}

func TestFindContainer(t *testing.T) {
	namespace := []string{"TestFindContainer"}
	containerID := "1"
	containerData := []types.Container{
		{ID: containerID},
	}
	containerJSONData := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID: containerID,
		},
	}

	dt := dockertest.New()
	dt.ProvideContainerList(containerData, nil)
	dt.ProvideContainerInspect(containerJSONData, nil)

	c, _ := New(ClientOption(dt.Client()))

	container, err := c.FindContainer(namespace)
	assert.Nil(t, err)
	assert.Equal(t, containerJSONData.ID, container.ID)

	assert.Equal(t, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "com.docker.stack.namespace=" + Namespace(namespace),
		}),
		Limit: 1,
	}, (<-dt.LastContainerList()).Options)

	assert.Equal(t, containerID, (<-dt.LastContainerInspect()).Container)
}

func TestNonExistentContainerStatus(t *testing.T) {
	namespace := []string{"namespace"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideContainerList(nil, dockertest.NotFoundErr{})

	status, err := c.Status(namespace)
	assert.Nil(t, err)
	assert.Equal(t, STOPPED, status)
}

func TestExistentContainerStatus(t *testing.T) {
	namespace := []string{"namespace"}
	containerID := "1"
	containerData := []types.Container{
		{ID: containerID},
	}
	containerJSONData := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID:    containerID,
			State: &types.ContainerState{},
		},
	}

	dt := dockertest.New()
	dt.ProvideContainerList(containerData, nil)
	dt.ProvideContainerInspect(containerJSONData, nil)

	c, _ := New(ClientOption(dt.Client()))
	status, err := c.Status(namespace)
	assert.Nil(t, err)
	assert.Equal(t, STOPPED, status)
}

func TestExistentContainerRunningStatus(t *testing.T) {
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
	dt.ProvideContainerList(containerData, nil)
	dt.ProvideContainerInspect(containerJSONData, nil)

	c, _ := New(ClientOption(dt.Client()))
	status, err := c.Status(namespace)
	assert.Nil(t, err)
	assert.Equal(t, RUNNING, status)
}
