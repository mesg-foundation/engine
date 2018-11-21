package container

import (
	"errors"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	dt := dockertest.New()
	c, err := New(ClientOption(dt.Client()))
	cfg, _ := config.Global()
	require.NoError(t, err)
	require.NotNil(t, c)

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

	require.Equal(t, "0.0.0.0:2377", (<-dt.LastSwarmInit()).Request.ListenAddr)

	ln := <-dt.LastNetworkCreate()
	require.Equal(t, cfg.Core.Name, ln.Name)
	require.Equal(t, types.NetworkCreate{
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
	require.NoError(t, err)
	require.NotNil(t, c)

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
	require.Equal(t, dockertest.NotFoundErr{}, err)

	require.Equal(t, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "com.docker.stack.namespace=" + c.Namespace(namespace),
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
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideContainerList(containerData, nil)
	dt.ProvideContainerInspect(containerJSONData, nil)

	container, err := c.FindContainer(namespace)
	require.NoError(t, err)
	require.Equal(t, containerJSONData.ID, container.ID)

	require.Equal(t, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "com.docker.stack.namespace=" + c.Namespace(namespace),
		}),
		Limit: 1,
	}, (<-dt.LastContainerList()).Options)

	require.Equal(t, containerID, (<-dt.LastContainerInspect()).Container)
}

func TestNonExistentContainerStatus(t *testing.T) {
	namespace := []string{"namespace"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideServiceInspectWithRaw(swarm.Service{}, nil, dockertest.NotFoundErr{})
	dt.ProvideContainerInspect(types.ContainerJSON{}, dockertest.NotFoundErr{})

	status, err := c.Status(namespace)
	require.NoError(t, err)
	require.Equal(t, STOPPED, status)

	resp := <-dt.LastServiceInspectWithRaw()
	require.Equal(t, c.Namespace(namespace), resp.ServiceID)
	require.Equal(t, types.ServiceInspectOptions{InsertDefaults: false}, resp.Options)
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
			State: &types.ContainerState{Running: true},
		},
	}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideServiceInspectWithRaw(swarm.Service{}, nil, nil)
	dt.ProvideContainerList(containerData, nil)
	dt.ProvideContainerInspect(containerJSONData, nil)

	status, err := c.Status(namespace)
	require.NoError(t, err)
	require.Equal(t, RUNNING, status)
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
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideContainerList(containerData, nil)
	dt.ProvideContainerInspect(containerJSONData, nil)

	status, err := c.Status(namespace)
	require.NoError(t, err)
	require.Equal(t, RUNNING, status)
}

func TestPresenceHandling(t *testing.T) {
	tests := []struct {
		param    error
		presence bool
		err      error
	}{
		{param: nil, presence: true, err: nil},
		{param: dockertest.NotFoundErr{}, presence: false, err: nil},
		{param: errors.New("test"), presence: false, err: errors.New("test")},
	}
	for _, test := range tests {
		presence, err := presenceHandling(test.param)
		require.Equal(t, test.presence, presence)
		require.Equal(t, test.err, err)
	}
}
