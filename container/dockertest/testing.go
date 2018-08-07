// Package dockertest provides functionalities for mocking and faking the Docker API.
package dockertest

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
)

// Testing provides functionalities to fake Docker API calls and responses.
type Testing struct {
	client *Client
}

// New creates a new Testing.
func New() *Testing {
	t := &Testing{
		client: newClient(),
	}
	return t
}

// Client returns a new mock client compatible with Docker's
// docker.CommonAPIClient interface.
func (t *Testing) Client() *Client {
	return t.client
}

// ProvideInfo sets fake return values for the next call to *Client.Into.
func (t *Testing) ProvideInfo(info types.Info, err error) {
	t.client.responses.info <- infoResponse{info, err}
}

// ProvideContainerList sets fake return values for the next call to *Client.ContainerList.
func (t *Testing) ProvideContainerList(containers []types.Container, err error) {
	t.client.responses.containerList <- containerListResponse{containers, err}
}

// ProvideContainerInspect sets fake return values for the next call to *Client.ContainerInspect.
func (t *Testing) ProvideContainerInspect(json types.ContainerJSON, err error) {
	t.client.responses.contanerInspect <- containerInspectResponse{json, err}
}

// ProvideImageBuild sets fake return values for the next call to *Client.ImageBuild.
func (t *Testing) ProvideImageBuild(rc io.ReadCloser, err error) {
	t.client.responses.imageBuild <- imageBuildResponse{types.ImageBuildResponse{Body: rc}, err}
}

// ProvideNetworkInspect sets fake return values for the next call to *Client.NetworkInspect.
func (t *Testing) ProvideNetworkInspect(resource types.NetworkResource, err error) {
	t.client.responses.networkInspect <- networkInspectResponse{resource, err}
}

// ProvideNetworkCreate sets fake return values for the next call to *Client.NetworkCreate.
func (t *Testing) ProvideNetworkCreate(response types.NetworkCreateResponse, err error) {
	t.client.responses.networkCreate <- networkCreateResponse{response, err}
}

// ProvideServiceList sets fake return values for the next call to *Client.ServiceList.
func (t *Testing) ProvideServiceList(services []swarm.Service, err error) {
	t.client.responses.serviceList <- serviceListResponse{services, err}
}

// ProvideServiceInspectWithRaw sets fake return values for the next call to *Client.ServiceInspectWithRaw.
func (t *Testing) ProvideServiceInspectWithRaw(service swarm.Service, data []byte, err error) {
	t.client.responses.serviceInspectWithRaw <- serviceInspectWithRawResponse{service, data, err}
}

// ProvideServiceLogs sets fake return values for the next call to *Client.ServiceLogs.
func (t *Testing) ProvideServiceLogs(rc io.ReadCloser, err error) {
	t.client.responses.serviceLogs <- serviceLogsResponse{rc, err}
}

// ProvideTaskList sets fake return values for the next call to *Client.TaskList.
func (t *Testing) ProvideTaskList(tasks []swarm.Task, err error) {
	t.client.responses.taskList <- taskListResponse{tasks, err}
}

// ProvideServiceCreate sets fake return values for the next call to *Client.ServiceCreate.
func (t *Testing) ProvideServiceCreate(response types.ServiceCreateResponse, err error) {
	// TODO(ilgooz) do the same shortcurt(calling needed Provides) made here
	// on other needed places too.
	t.ProvideContainerList([]types.Container{}, nil)
	t.ProvideContainerInspect(types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			State: &types.ContainerState{Running: true},
		},
	}, nil)

	t.client.responses.serviceCreate <- serviceCreateResponse{response, err}
}

// ProvideServiceRemove sets fake return values for the next call to *Client.ServiceRemove.
func (t *Testing) ProvideServiceRemove(err error) {
	t.client.responses.serviceRemove <- serviceRemoveResponse{err}
}

// LastContainerList returns a channel that receives call arguments of last *Client.ContainerList call.
func (t *Testing) LastContainerList() <-chan ContainerListRequest {
	return t.client.requests.containerList
}

// LastContainerInspect returns a channel that receives call arguments of last *Client.ContainerInspect call.
func (t *Testing) LastContainerInspect() <-chan ContainerInspectRequest {
	return t.client.requests.containerInspect
}

// LastSwarmInit returns a channel that receives call arguments of last *Client.SwarmInit call.
func (t *Testing) LastSwarmInit() <-chan SwarmInitRequest {
	return t.client.requests.swarmInit
}

// LastNetworkCreate returns a channel that receives call arguments of last *Client.NetworkCreate call.
func (t *Testing) LastNetworkCreate() <-chan NetworkCreateRequest {
	return t.client.requests.networkCreate
}

// LastNetworkInspect returns a channel that receives call arguments of last *Client.NetworkInspect call.
func (t *Testing) LastNetworkInspect() <-chan NetworkInspectRequest {
	return t.client.requests.networkInspect
}

// LastInfo returns a channel that receives call arguments of last *Client.Info call.
func (t *Testing) LastInfo() <-chan InfoRequest {
	return t.client.requests.info
}

// LastNegotiateAPIVersion returns a channel that receives call arguments of last *Client.NegotiateAPIVersion call.
func (t *Testing) LastNegotiateAPIVersion() <-chan NegotiateAPIVersionRequest {
	return t.client.requests.negotiateAPIVersion
}

// LastImageBuild returns a channel that receives call arguments of last *Client.ImageBuild call.
func (t *Testing) LastImageBuild() <-chan ImageBuildRequest {
	return t.client.requests.imageBuild
}

// LastNetworkRemove returns a channel that receives call arguments of last *Client.NetworkRemove call.
func (t *Testing) LastNetworkRemove() <-chan NetworkRemoveRequest {
	return t.client.requests.networkRemove
}

// LastTaskList returns a channel that receives call arguments of last *Client.TaskList call.
func (t *Testing) LastTaskList() <-chan TaskListRequest {
	return t.client.requests.taskList
}

// LastServiceCreate returns a channel that receives call arguments of last *Client.ServiceCreate call.
func (t *Testing) LastServiceCreate() <-chan ServiceCreateRequest {
	return t.client.requests.serviceCreate
}

// LastServiceList returns a channel that receives call arguments of last *Client.ServiceList call.
func (t *Testing) LastServiceList() <-chan ServiceListRequest {
	return t.client.requests.serviceList
}

// LastServiceInspectWithRaw returns a channel that receives call arguments of last *Client.ServiceInspectWithRaw call.
func (t *Testing) LastServiceInspectWithRaw() <-chan ServiceInspectWithRawRequest {
	return t.client.requests.serviceInspectWithRaw
}

// LastServiceRemove returns a channel that receives call arguments of last *Client.ServiceRemove call.
func (t *Testing) LastServiceRemove() <-chan ServiceRemoveRequest {
	return t.client.requests.serviceRemove
}

// LastServiceLogs returns a channel that receives call arguments of last *Client.ServiceLogs call.
func (t *Testing) LastServiceLogs() <-chan ServiceLogsRequest {
	return t.client.requests.serviceLogs
}
