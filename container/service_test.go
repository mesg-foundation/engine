package container

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStartService(t *testing.T) {
	var (
		namespace = []string{"1"}
		serviceID = "2"
		options   = ServiceOptions{
			Image:     "3",
			Namespace: namespace,
		}
		c, mc         = newTesting(t)
		fullNamespace = c.Namespace(namespace)
	)

	// status:
	mockStatus(t, mc, fullNamespace, STOPPED)

	// create service:
	mc.On("ServiceCreate", mock.Anything, options.toSwarmServiceSpec(c), types.ServiceCreateOptions{}).
		Once().
		Return(types.ServiceCreateResponse{ID: serviceID}, nil)

	// wait status:
	mockWaitForStatus(t, mc, fullNamespace, RUNNING)

	id, err := c.StartService(options)
	require.NoError(t, err)
	require.Equal(t, serviceID, id)

	mc.AssertExpectations(t)
}

func TestStartServiceRunningStatus(t *testing.T) {
	var (
		namespace = []string{"1"}
		serviceID = "2"
		options   = ServiceOptions{
			Image:     "3",
			Namespace: namespace,
		}
		c, mc         = newTesting(t)
		fullNamespace = c.Namespace(namespace)
	)

	// status:
	mockStatus(t, mc, fullNamespace, RUNNING)

	// find service:
	mc.On("ServiceInspectWithRaw", mock.Anything, mock.Anything, mock.Anything).
		Once().
		Return(swarm.Service{ID: serviceID}, nil, nil)

	id, err := c.StartService(options)
	require.NoError(t, err)
	require.Equal(t, serviceID, id)

	mc.AssertExpectations(t)
}

func TestStopService(t *testing.T) {
	var (
		c, mc         = newTesting(t)
		containerID   = "1"
		namespace     = []string{"2"}
		fullNamespace = c.Namespace(namespace)
	)

	mockStatus(t, mc, fullNamespace, RUNNING)

	// remove service:
	mc.On("ServiceRemove", mock.Anything, fullNamespace).
		Once().
		Return(nil)

	// delete pending container:
	mc.On("ContainerList", mock.Anything, mock.Anything).
		Once().
		Return([]types.Container{{ID: containerID}}, nil)

	mc.On("ContainerInspect", mock.Anything, mock.Anything).
		Once().
		Return(types.ContainerJSON{
			ContainerJSONBase: &types.ContainerJSONBase{ID: containerID},
		}, nil)

	mc.On("ContainerStop", mock.Anything, containerID, mock.AnythingOfType("*time.Duration")).
		Once().
		Return(nil)

	mc.On("ContainerRemove", mock.Anything, containerID, types.ContainerRemoveOptions{}).
		Once().
		Return(nil)

	mc.On("ContainerList", mock.Anything, mock.Anything).
		Once().
		Return(nil, dockertest.NotFoundErr{})

	// wait status:
	mockWaitForStatus(t, mc, fullNamespace, STOPPED)

	require.NoError(t, c.StopService(namespace))

	mc.AssertExpectations(t)
}

func TestStopNotExistingService(t *testing.T) {
	namespace := []string{"namespace"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideContainerList([]types.Container{}, nil)
	dt.ProvideContainerInspect(types.ContainerJSON{}, dockertest.NotFoundErr{})
	dt.ProvideServiceRemove(dockertest.NotFoundErr{})
	dt.ProvideServiceInspectWithRaw(swarm.Service{}, nil, dockertest.NotFoundErr{})
	dt.ProvideContainerInspect(types.ContainerJSON{}, dockertest.NotFoundErr{})

	require.NoError(t, c.StopService(namespace))
}

func TestFindService(t *testing.T) {
	namespace := []string{"namespace"}
	swarmService := swarm.Service{ID: "1"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideServiceInspectWithRaw(swarmService, nil, nil)

	service, err := c.FindService(namespace)
	require.NoError(t, err)
	require.Equal(t, swarmService.ID, service.ID)

	li := <-dt.LastServiceInspectWithRaw()
	require.Equal(t, c.Namespace(namespace), li.ServiceID)
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
	require.Equal(t, c.Namespace(namespace), li.ServiceID)
	require.Equal(t, types.ServiceInspectOptions{}, li.Options)
}

func TestListServices(t *testing.T) {
	namespace := []string{"namespace"}
	namespace1 := []string{"namespace"}
	label := "1"
	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))
	swarmServices := []swarm.Service{
		{Spec: swarm.ServiceSpec{Annotations: swarm.Annotations{Name: c.Namespace(namespace)}}},
		{Spec: swarm.ServiceSpec{Annotations: swarm.Annotations{Name: c.Namespace(namespace1)}}},
	}

	dt.ProvideServiceList(swarmServices, nil)

	services, err := c.ListServices(label)
	require.NoError(t, err)
	require.Equal(t, 2, len(services))
	require.Equal(t, c.Namespace(namespace), services[0].Spec.Name)
	require.Equal(t, c.Namespace(namespace1), services[1].Spec.Name)

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
	require.NoError(t, err)
	defer reader.Close()

	bytes, err := ioutil.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, data, bytes)

	ll := <-dt.LastServiceLogs()
	require.Equal(t, c.Namespace(namespace), ll.ServiceID)
	require.Equal(t, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: false,
		Follow:     true,
	}, ll.Options)
}
