package dockertest

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
)

// NetworkCreateRequest holds call arguments of *Client.NetworkCreate.
type NetworkCreateRequest struct {
	Name    string
	Options types.NetworkCreate
}

// ImageBuildRequest holds call arguments of *Client.ImageBuild.
type ImageBuildRequest struct {
	FileData []byte
	Options  types.ImageBuildOptions
}

// NetworkInspectRequest holds call arguments of *Client.NetworkInspect.
type NetworkInspectRequest struct {
	Network string
	Options types.NetworkInspectOptions
}

// ServiceCreateRequest holds call arguments of *Client.ServiceCreate.
type ServiceCreateRequest struct {
	Service swarm.ServiceSpec
	Options types.ServiceCreateOptions
}

// ServiceInspectWithRawRequest holds call arguments of *Client.ServiceInspectWithRaw.
type ServiceInspectWithRawRequest struct {
	ServiceID string
	Options   types.ServiceInspectOptions
}

// ServiceLogsRequest holds call arguments of *Client.ServiceLogs.
type ServiceLogsRequest struct {
	ServiceID string
	Options   types.ContainerLogsOptions
}

// NegotiateAPIVersionRequest holds call arguments of *Client.NegotiateAPIVersion.
type NegotiateAPIVersionRequest struct {
}

// SwarmInitRequest holds call arguments of *Client.SwarmInit.
type SwarmInitRequest struct {
	Request swarm.InitRequest
}

// InfoRequest holds call arguments of *Client.Info.
type InfoRequest struct {
}

// ContainerListRequest holds call arguments of *Client.ContainerList.
type ContainerListRequest struct {
	Options types.ContainerListOptions
}

// ContainerInspectRequest holds call arguments of *Client.ContainerInspect.
type ContainerInspectRequest struct {
	Container string
}

// ContainerStopRequest holds call arguments of *Client.ContainerStop.
type ContainerStopRequest struct {
	Container string
}

// ContainerRemoveRequest holds call arguments of *Client.ContainerRemove.
type ContainerRemoveRequest struct {
	Container string
	Options   types.ContainerRemoveOptions
}

// NetworkRemoveRequest holds call arguments of *Client.NetworkRemove.
type NetworkRemoveRequest struct {
	Network string
}

// TaskListRequest holds call arguments of *Client.TaskList.
type TaskListRequest struct {
	Options types.TaskListOptions
}

// ServiceListRequest holds call arguments of *Client.ServiceList.
type ServiceListRequest struct {
	Options types.ServiceListOptions
}

// ServiceRemoveRequest holds call arguments of *Client.ServiceRemove.
type ServiceRemoveRequest struct {
	ServiceID string
}
