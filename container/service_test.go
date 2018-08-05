package container

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/stvp/assert"
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
	assert.Nil(t, err)
	assert.Equal(t, containerID, id)

	ls := <-dt.LastServiceCreate()
	assert.Equal(t, options.toSwarmServiceSpec(), ls.Service)
	assert.Equal(t, types.ServiceCreateOptions{}, ls.Options)
}

func TestStopService(t *testing.T) {
	namespace := []string{"namespace"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	assert.Nil(t, c.StopService(namespace))

	li := <-dt.LastServiceInspectWithRaw()
	assert.Equal(t, Namespace(namespace), li.ServiceID)
	assert.Equal(t, types.ServiceInspectOptions{}, li.Options)

	assert.Equal(t, Namespace(namespace), <-dt.LastServiceRemove())
}

func TestStopNotExistingService(t *testing.T) {
	namespace := []string{"namespace"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideServiceInspectWithRaw(swarm.Service{}, nil, dockertest.NotFoundErr{})

	assert.Nil(t, c.StopService(namespace))

	li := <-dt.LastServiceInspectWithRaw()
	assert.Equal(t, Namespace(namespace), li.ServiceID)
	assert.Equal(t, types.ServiceInspectOptions{}, li.Options)

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
	assert.Nil(t, err)
	assert.Equal(t, STOPPED, status)

	li := <-dt.LastServiceInspectWithRaw()
	assert.Equal(t, Namespace(namespace), li.ServiceID)
	assert.Equal(t, types.ServiceInspectOptions{}, li.Options)
}

func TestIntegrationServiceStatusRunning(t *testing.T) {
	namespace := []string{"namespace"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	status, err := c.ServiceStatus(namespace)
	assert.Nil(t, err)
	assert.Equal(t, RUNNING, status)

	li := <-dt.LastServiceInspectWithRaw()
	assert.Equal(t, Namespace(namespace), li.ServiceID)
	assert.Equal(t, types.ServiceInspectOptions{}, li.Options)
}

func TestFindService(t *testing.T) {
	namespace := []string{"namespace"}
	swarmService := swarm.Service{ID: "1"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideServiceInspectWithRaw(swarmService, nil, nil)

	service, err := c.FindService(namespace)
	assert.Nil(t, err)
	assert.Equal(t, swarmService.ID, service.ID)

	li := <-dt.LastServiceInspectWithRaw()
	assert.Equal(t, Namespace(namespace), li.ServiceID)
	assert.Equal(t, types.ServiceInspectOptions{}, li.Options)
}

func TestFindServiceNotExisting(t *testing.T) {
	namespace := []string{"namespace"}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideServiceInspectWithRaw(swarm.Service{}, nil, dockertest.NotFoundErr{})

	_, err := c.FindService(namespace)
	assert.Equal(t, dockertest.NotFoundErr{}, err)

	li := <-dt.LastServiceInspectWithRaw()
	assert.Equal(t, Namespace(namespace), li.ServiceID)
	assert.Equal(t, types.ServiceInspectOptions{}, li.Options)
}

func TestListServices(t *testing.T) {
	namespace := []string{"namespace"}
	label := "1"
	swarmServices := []swarm.Service{
		{Spec: swarm.ServiceSpec{Annotations: swarm.Annotations{Name: Namespace(namespace)}}},
		{Spec: swarm.ServiceSpec{Annotations: swarm.Annotations{Name: Namespace(namespace)}}},
	}

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideServiceList(swarmServices, nil)

	services, err := c.ListServices(label)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(services))
	assert.Equal(t, Namespace(namespace), services[0].Spec.Name)
	assert.Equal(t, Namespace(namespace), services[1].Spec.Name)

	assert.Equal(t, types.ServiceListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: label,
		}),
	}, <-dt.LastServiceList())
}

func TestServiceLogs(t *testing.T) {
	namespace := []string{"namespace"}
	data := "mesg"

	dt := dockertest.New()
	c, _ := New(ClientOption(dt.Client()))

	dt.ProvideServiceLogs(ioutil.NopCloser(strings.NewReader(data)), nil)

	reader, err := c.ServiceLogs(namespace)
	assert.Nil(t, err)

	defer reader.Close()
	bytes, err := ioutil.ReadAll(reader)
	assert.Nil(t, err)
	assert.Equal(t, data, string(bytes))

	ll := <-dt.LastServiceLogs()

	assert.Equal(t, Namespace(namespace), ll.ServiceID)
	assert.Equal(t, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: false,
		Follow:     true,
	}, ll.Options)
}
