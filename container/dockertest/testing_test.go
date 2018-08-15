package dockertest

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/stretchr/testify/require"
)

var errGeneric = errors.New("titan of the errors")

func TestNew(t *testing.T) {
	dt := New()
	require.NotNil(t, dt)
}

func TestClient(t *testing.T) {
	dt := New()
	require.NotNil(t, dt.Client())
}

func TestNegotiateAPIVersion(t *testing.T) {
	dt := New()
	dt.Client().NegotiateAPIVersion(context.Background())

	select {
	case <-dt.LastNegotiateAPIVersion():
	default:
		t.Error("last call to *Client.NegotiateAPIVersion should be observable")
	}
}

func TestNetworkInspect(t *testing.T) {
	resource := types.NetworkResource{ID: "1"}
	network := "2"
	options := types.NetworkInspectOptions{Verbose: true}

	dt := New()
	dt.ProvideNetworkInspect(resource, errGeneric)

	resource1, err := dt.Client().NetworkInspect(context.Background(), network, options)
	require.Equal(t, errGeneric, err)
	require.Equal(t, resource, resource1)

	ll := <-dt.LastNetworkInspect()
	require.Equal(t, network, ll.Network)
	require.Equal(t, options, ll.Options)
}

func TestNetworkCreate(t *testing.T) {
	response := types.NetworkCreateResponse{ID: "1"}
	name := "2"
	options := types.NetworkCreate{CheckDuplicate: true}

	dt := New()
	dt.ProvideNetworkCreate(response, errGeneric)

	response1, err := dt.Client().NetworkCreate(context.Background(), name, options)
	require.Equal(t, errGeneric, err)
	require.Equal(t, response, response1)

	ll := <-dt.LastNetworkCreate()
	require.Equal(t, name, ll.Name)
	require.Equal(t, options, ll.Options)
}

func TestSwarmInit(t *testing.T) {
	request := swarm.InitRequest{ForceNewCluster: true}

	dt := New()

	data, err := dt.Client().SwarmInit(context.Background(), request)
	require.Nil(t, err)
	require.Equal(t, "", data)

	require.Equal(t, request, (<-dt.LastSwarmInit()).Request)
}

func TestInfo(t *testing.T) {
	info := types.Info{ID: "1"}

	dt := New()
	dt.ProvideInfo(info, errGeneric)

	info1, err := dt.Client().Info(context.Background())
	require.Equal(t, errGeneric, err)
	require.Equal(t, info, info1)

	select {
	case <-dt.LastInfo():
	default:
		t.Error("last call to *Client.Info should be observable")
	}
}

func TestContainerList(t *testing.T) {
	containers := []types.Container{{ID: "1"}, {ID: "2"}}
	options := types.ContainerListOptions{Quiet: true}

	dt := New()
	dt.ProvideContainerList(containers, errGeneric)

	containers1, err := dt.Client().ContainerList(context.Background(), options)
	require.Equal(t, errGeneric, err)
	require.Equal(t, containers, containers1)

	ll := <-dt.LastContainerList()
	require.Equal(t, options, ll.Options)
}

func TestContainerInspect(t *testing.T) {
	container := "1"
	containerJSON := types.ContainerJSON{ContainerJSONBase: &types.ContainerJSONBase{ID: "2"}}

	dt := New()
	dt.ProvideContainerInspect(containerJSON, errGeneric)

	containerJSON1, err := dt.Client().ContainerInspect(context.Background(), container)
	require.Equal(t, errGeneric, err)
	require.Equal(t, containerJSON, containerJSON1)

	ll := <-dt.LastContainerInspect()
	require.Equal(t, container, ll.Container)
}

func TestImageBuild(t *testing.T) {
	options := types.ImageBuildOptions{SuppressOutput: true}
	response := []byte{1, 2}
	request := []byte{3}

	dt := New()
	dt.ProvideImageBuild(ioutil.NopCloser(bytes.NewReader(response)), errGeneric)

	resp, err := dt.Client().ImageBuild(context.Background(),
		ioutil.NopCloser(bytes.NewReader(request)), options)
	require.Equal(t, errGeneric, err)
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	require.Nil(t, err)
	require.Equal(t, response, respData)

	ll := <-dt.LastImageBuild()
	require.Equal(t, options, ll.Options)
	require.Equal(t, request, ll.FileData)
}

func TestNetworkRemove(t *testing.T) {
	network := "1"
	dt := New()
	require.Nil(t, dt.Client().NetworkRemove(context.Background(), network))
	require.Equal(t, network, (<-dt.LastNetworkRemove()).Network)
}

func TestTaskList(t *testing.T) {
	tasks := []swarm.Task{{ID: "1"}, {ID: "2"}}
	options := types.TaskListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key: "1",
		}),
	}

	dt := New()
	dt.ProvideTaskList(tasks, errGeneric)

	tasks1, err := dt.Client().TaskList(context.Background(), options)
	require.Equal(t, errGeneric, err)
	require.Equal(t, tasks, tasks1)

	require.Equal(t, options, (<-dt.LastTaskList()).Options)
}

func TestServiceCreate(t *testing.T) {
	response := types.ServiceCreateResponse{ID: "1"}
	service := swarm.ServiceSpec{Annotations: swarm.Annotations{Name: "1"}}
	options := types.ServiceCreateOptions{QueryRegistry: true}

	dt := New()
	dt.ProvideServiceCreate(response, errGeneric)

	response1, err := dt.Client().ServiceCreate(context.Background(), service, options)
	require.Equal(t, errGeneric, err)
	require.Equal(t, response, response1)

	ll := <-dt.LastServiceCreate()
	require.Equal(t, service, ll.Service)
	require.Equal(t, options, ll.Options)
}

func TestServiceList(t *testing.T) {
	services := []swarm.Service{{ID: "1"}, {ID: "2"}}
	options := types.ServiceListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key: "label",
		}),
	}

	dt := New()
	dt.ProvideServiceList(services, errGeneric)

	services1, err := dt.Client().ServiceList(context.Background(), options)
	require.Equal(t, errGeneric, err)
	require.Equal(t, services, services1)

	ll := <-dt.LastServiceList()
	require.Equal(t, options, ll.Options)
}

func TestServiceInspectWithRaw(t *testing.T) {
	serviceID := "1"
	options := types.ServiceInspectOptions{InsertDefaults: true}
	service := swarm.Service{ID: "1"}
	data := []byte{1, 2}

	dt := New()
	dt.ProvideServiceInspectWithRaw(service, data, errGeneric)

	service1, data1, err := dt.Client().ServiceInspectWithRaw(context.Background(), serviceID, options)
	require.Equal(t, errGeneric, err)
	require.Equal(t, service, service1)
	require.Equal(t, data, data1)

	ll := <-dt.LastServiceInspectWithRaw()
	require.Equal(t, serviceID, ll.ServiceID)
	require.Equal(t, options, ll.Options)
}

func TestServiceRemove(t *testing.T) {
	serviceID := "1"

	dt := New()
	dt.ProvideServiceRemove(errGeneric)

	require.Equal(t, errGeneric, dt.Client().ServiceRemove(context.Background(), serviceID))

	ll := <-dt.LastServiceRemove()
	require.Equal(t, serviceID, ll.ServiceID)
}

func TestServiceLogs(t *testing.T) {
	serviceID := "1"
	data := []byte{1, 2}
	options := types.ContainerLogsOptions{ShowStdout: true}

	dt := New()
	dt.ProvideServiceLogs(ioutil.NopCloser(bytes.NewReader(data)), errGeneric)

	rc, err := dt.Client().ServiceLogs(context.Background(), serviceID, options)
	require.Equal(t, errGeneric, err)
	defer rc.Close()

	data1, err := ioutil.ReadAll(rc)
	require.Nil(t, err)
	require.Equal(t, data, data1)

	ll := <-dt.LastServiceLogs()
	require.Equal(t, serviceID, ll.ServiceID)
	require.Equal(t, options, ll.Options)
}
