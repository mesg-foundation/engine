package container

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/stretchr/testify/require"
)

func TestStartService(t *testing.T) {
	namespace := []string{"namespace"}
	containerID := "id"
	options := ServiceOptions{
		Image:     "nginx",
		Namespace: namespace,
	}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideServiceCreate(types.ServiceCreateResponse{ID: containerID}, nil)

	id, err := c.StartService(options)
	require.Nil(t, err)
	require.Equal(t, containerID, id)

	ls := <-dt.LastServiceCreate()
	require.Equal(t, options.toSwarmServiceSpec(), ls.Service)
	require.Equal(t, types.ServiceCreateOptions{}, ls.Options)
}

func TestStopService(t *testing.T) {
	namespace := []string{"namespace"}
	containerData := []types.Container{}
	containerJSONData := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			State: &types.ContainerState{},
		},
	}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideContainerList(containerData, nil)
	dt.ProvideContainerInspect(containerJSONData, nil)

	require.Nil(t, c.StopService(namespace))

	li := <-dt.LastServiceInspectWithRaw()
	require.Equal(t, Namespace(namespace), li.ServiceID)
	require.Equal(t, types.ServiceInspectOptions{}, li.Options)

	require.Equal(t, Namespace(namespace), (<-dt.LastServiceRemove()).ServiceID)
}

func TestStopNotExistingService(t *testing.T) {
	namespace := []string{"namespace"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideServiceInspectWithRaw(swarm.Service{}, nil, dockertest.NotFoundErr{})

	require.Nil(t, c.StopService(namespace))

	li := <-dt.LastServiceInspectWithRaw()
	require.Equal(t, Namespace(namespace), li.ServiceID)
	require.Equal(t, types.ServiceInspectOptions{}, li.Options)

	select {
	case <-dt.LastServiceRemove():
		t.Error("should not remove non existent service")
	default:
	}
}

func TestServiceStatusNeverStarted(t *testing.T) {
	namespace := []string{"namespace"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideServiceInspectWithRaw(swarm.Service{}, nil, dockertest.NotFoundErr{})

	status, err := c.ServiceStatus(namespace)
	require.Nil(t, err)
	require.Equal(t, STOPPED, status)

	li := <-dt.LastServiceInspectWithRaw()
	require.Equal(t, Namespace(namespace), li.ServiceID)
	require.Equal(t, types.ServiceInspectOptions{}, li.Options)
}

func TestServiceStatusRunning(t *testing.T) {
	namespace := []string{"namespace"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	status, err := c.ServiceStatus(namespace)
	require.Nil(t, err)
	require.Equal(t, RUNNING, status)

	li := <-dt.LastServiceInspectWithRaw()
	require.Equal(t, Namespace(namespace), li.ServiceID)
	require.Equal(t, types.ServiceInspectOptions{}, li.Options)
}

func TestFindService(t *testing.T) {
	namespace := []string{"namespace"}
	swarmService := swarm.Service{ID: "1"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideServiceInspectWithRaw(swarmService, nil, nil)

	service, err := c.FindService(namespace)
	require.Nil(t, err)
	require.Equal(t, swarmService.ID, service.ID)

	li := <-dt.LastServiceInspectWithRaw()
	require.Equal(t, Namespace(namespace), li.ServiceID)
	require.Equal(t, types.ServiceInspectOptions{}, li.Options)
}

func TestFindServiceNotExisting(t *testing.T) {
	namespace := []string{"namespace"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideServiceInspectWithRaw(swarm.Service{}, nil, dockertest.NotFoundErr{})

	_, err := c.FindService(namespace)
	require.Equal(t, dockertest.NotFoundErr{}, err)

	li := <-dt.LastServiceInspectWithRaw()
	require.Equal(t, Namespace(namespace), li.ServiceID)
	require.Equal(t, types.ServiceInspectOptions{}, li.Options)
}

func TestListServices(t *testing.T) {
	namespace := []string{"namespace"}
	namespace1 := []string{"namespace"}
	label := "1"
	swarmServices := []swarm.Service{
		{Spec: swarm.ServiceSpec{Annotations: swarm.Annotations{Name: Namespace(namespace)}}},
		{Spec: swarm.ServiceSpec{Annotations: swarm.Annotations{Name: Namespace(namespace1)}}},
	}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideServiceList(swarmServices, nil)

	services, err := c.ListServices(label)
	require.Nil(t, err)
	require.Equal(t, 2, len(services))
	require.Equal(t, Namespace(namespace), services[0].Spec.Name)
	require.Equal(t, Namespace(namespace1), services[1].Spec.Name)

	require.Equal(t, types.ServiceListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: label,
		}),
	}, (<-dt.LastServiceList()).Options)
}

func TestServiceLogs(t *testing.T) {
	namespace := []string{"namespace"}
	data := []byte{1, 2}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideServiceLogs(ioutil.NopCloser(bytes.NewReader(data)), nil)

	reader, err := c.ServiceLogs(namespace)
	require.Nil(t, err)
	defer reader.Close()

	bytes, err := ioutil.ReadAll(reader)
	require.Nil(t, err)
	require.Equal(t, data, bytes)

	ll := <-dt.LastServiceLogs()
	require.Equal(t, Namespace(namespace), ll.ServiceID)
	require.Equal(t, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: false,
		Follow:     true,
	}, ll.Options)
}
