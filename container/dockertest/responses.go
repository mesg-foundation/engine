package dockertest

import (
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/swarm"
)

// serviceLogsResponse holds fake return values of *Client.ServiceLogs.
type serviceLogsResponse struct {
	rc  io.ReadCloser
	err error
}

// serviceRemoveResponse holds fake return values of *Client.ServiceRemove.
type serviceRemoveResponse struct {
	err error
}

// serviceInspectWithRawResponse holds fake return values of *Client.ServiceInspectWithRaw.
type serviceInspectWithRawResponse struct {
	service swarm.Service
	data    []byte
	err     error
}

// serviceListResponse holds fake return values of *Client.ServiceList.
type serviceListResponse struct {
	services []swarm.Service
	err      error
}

// serviceCreateResponse holds fake return values of *Client.ServiceCreate.
type serviceCreateResponse struct {
	response types.ServiceCreateResponse
	err      error
}

// taskListResponse holds fake return values of *Client.TaskList.
type taskListResponse struct {
	tasks []swarm.Task
	err   error
}

// networkCreateResponse holds fake return values of *Client.NetworkCreate.
type networkCreateResponse struct {
	response types.NetworkCreateResponse
	err      error
}

// networkInspectResponse holds fake return values of *Client.NetworkInspect.
type networkInspectResponse struct {
	resource types.NetworkResource
	err      error
}

// imageBuildResponse holds fake return values of *Client.ImageBuild.
type imageBuildResponse struct {
	response types.ImageBuildResponse
	err      error
}

// infoResponse holds fake return values of *Client.Info.
type infoResponse struct {
	info types.Info
	err  error
}

// containerListResponse holds fake return values of *Client.ContainerList.
type containerListResponse struct {
	containers []types.Container
	err        error
}

// containerInspectResponse holds fake return values of *Client.ContainerInspect.
type containerInspectResponse struct {
	json types.ContainerJSON
	err  error
}

// containerStopResponse holds fake return values of *Client.ContainerStop.
type containerStopResponse struct {
	err error
}

// containerRemoveResponse holds fake return values of *Client.ContainerRemove.
type containerRemoveResponse struct {
	err error
}

// eventsResponse holds fake return values of *Client.Events.
type eventsResponse struct {
	messageChan <-chan events.Message
	errChan     <-chan error
}
