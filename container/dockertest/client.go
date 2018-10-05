package dockertest

import (
	"context"
	"io"
	"io/ioutil"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/docker/docker/client"
)

// Client satisfies docker.CommonAPIClient for mocking.
type Client struct {
	requests  *requests
	responses *responses

	// a placeholder for unimplemented methods.
	docker.CommonAPIClient
}

// requests encapsulates request channels that keeps call arguments of Docker client methods.
// Each call to a Docker client method piped to its request channel so multiple
// call to a same method can be inspected by reading its channel multiple times.
// Inspecting call arguments can be made by listening chan returned by LastX methods of *Testing.
type requests struct {
	negotiateAPIVersion   chan NegotiateAPIVersionRequest
	swarmInit             chan SwarmInitRequest
	networkCreate         chan NetworkCreateRequest
	info                  chan InfoRequest
	containerList         chan ContainerListRequest
	containerInspect      chan ContainerInspectRequest
	containerStop         chan ContainerStopRequest
	containerRemove       chan ContainerRemoveRequest
	imageBuild            chan ImageBuildRequest
	networkInspect        chan NetworkInspectRequest
	networkRemove         chan NetworkRemoveRequest
	taskList              chan TaskListRequest
	serviceCreate         chan ServiceCreateRequest
	serviceList           chan ServiceListRequest
	serviceInspectWithRaw chan ServiceInspectWithRawRequest
	serviceRemove         chan ServiceRemoveRequest
	serviceLogs           chan ServiceLogsRequest
	events                chan EventsRequest
}

// responses encapsulates response channels that holds 'faked' return values of Docker client methods.
// To fake a Docker client method's response send fake return values to it's
// response channel. This can be made by calling ProvideX methods of *Testing.
// We use channels here instead of setting values on the struct in case of a need for returning
// conditional responses depending on request paramaters in future to deal with parallel calls made
// to same client methods.
type responses struct {
	info                  chan infoResponse
	imageBuild            chan imageBuildResponse
	networkInspect        chan networkInspectResponse
	networkCreate         chan networkCreateResponse
	taskList              chan taskListResponse
	serviceCreate         chan serviceCreateResponse
	serviceList           chan serviceListResponse
	serviceInspectWithRaw chan serviceInspectWithRawResponse
	serviceRemove         chan serviceRemoveResponse
	serviceLogs           chan serviceLogsResponse
	containerInspect      chan containerInspectResponse
	containerList         chan containerListResponse
	containerStop         chan containerStopResponse
	containerRemove       chan containerRemoveResponse
	events                chan eventsResponse
}

// newClient returns a new mock Client for Docker.
func newClient() *Client {
	return &Client{
		// buffered channels helps with writing tests in synchronous syntax while using
		// LastX methods of Testing.
		requests: &requests{
			negotiateAPIVersion:   make(chan NegotiateAPIVersionRequest, 20),
			swarmInit:             make(chan SwarmInitRequest, 20),
			networkCreate:         make(chan NetworkCreateRequest, 20),
			info:                  make(chan InfoRequest, 20),
			containerList:         make(chan ContainerListRequest, 20),
			containerInspect:      make(chan ContainerInspectRequest, 20),
			containerStop:         make(chan ContainerStopRequest, 20),
			containerRemove:       make(chan ContainerRemoveRequest, 20),
			imageBuild:            make(chan ImageBuildRequest, 20),
			networkInspect:        make(chan NetworkInspectRequest, 20),
			networkRemove:         make(chan NetworkRemoveRequest, 20),
			taskList:              make(chan TaskListRequest, 20),
			serviceCreate:         make(chan ServiceCreateRequest, 20),
			serviceList:           make(chan ServiceListRequest, 20),
			serviceInspectWithRaw: make(chan ServiceInspectWithRawRequest, 20),
			serviceRemove:         make(chan ServiceRemoveRequest, 20),
			serviceLogs:           make(chan ServiceLogsRequest, 20),
			events:                make(chan EventsRequest, 20),
		},

		// buffered channels helps with writing tests in synchronous syntax while using
		// ProvideX methods of Testing.
		responses: &responses{
			info:                  make(chan infoResponse, 20),
			imageBuild:            make(chan imageBuildResponse, 20),
			networkInspect:        make(chan networkInspectResponse, 20),
			networkCreate:         make(chan networkCreateResponse, 20),
			taskList:              make(chan taskListResponse, 20),
			serviceCreate:         make(chan serviceCreateResponse, 20),
			serviceList:           make(chan serviceListResponse, 20),
			serviceInspectWithRaw: make(chan serviceInspectWithRawResponse, 20),
			serviceRemove:         make(chan serviceRemoveResponse, 20),
			serviceLogs:           make(chan serviceLogsResponse, 20),
			containerInspect:      make(chan containerInspectResponse, 20),
			containerList:         make(chan containerListResponse, 20),
			containerStop:         make(chan containerStopResponse, 20),
			containerRemove:       make(chan containerRemoveResponse, 20),
			events:                make(chan eventsResponse, 20),
		},
	}
}

// note that select: default statements are a short cut of returning zero value responses
// instead of an explicit need to call Provide methods with zero values.

// NegotiateAPIVersion is the mock version of the actual method.
func (c *Client) NegotiateAPIVersion(ctx context.Context) {
	c.requests.negotiateAPIVersion <- NegotiateAPIVersionRequest{}
}

// NetworkInspect is the mock version of the actual method.
func (c *Client) NetworkInspect(ctx context.Context, network string,
	options types.NetworkInspectOptions) (types.NetworkResource, error) {
	c.requests.networkInspect <- NetworkInspectRequest{network, options}
	select {
	case resp := <-c.responses.networkInspect:
		return resp.resource, resp.err
	default:
		return types.NetworkResource{}, nil
	}
}

// NetworkCreate is the mock version of the actual method.
func (c *Client) NetworkCreate(ctx context.Context, name string,
	options types.NetworkCreate) (types.NetworkCreateResponse, error) {
	c.requests.networkCreate <- NetworkCreateRequest{name, options}
	select {
	case resp := <-c.responses.networkCreate:
		return resp.response, resp.err
	default:
		return types.NetworkCreateResponse{}, nil
	}
}

// SwarmInit is the mock version of the actual method.
func (c *Client) SwarmInit(ctx context.Context, req swarm.InitRequest) (string, error) {
	c.requests.swarmInit <- SwarmInitRequest{req}
	return "", nil
}

// Info is the mock version of the actual method.
func (c *Client) Info(context.Context) (types.Info, error) {
	c.requests.info <- InfoRequest{}
	select {
	case resp := <-c.responses.info:
		return resp.info, resp.err
	default:
		return types.Info{}, nil
	}
}

// ContainerList is the mock version of the actual method.
func (c *Client) ContainerList(ctx context.Context,
	options types.ContainerListOptions) ([]types.Container, error) {
	c.requests.containerList <- ContainerListRequest{options}
	select {
	case resp := <-c.responses.containerList:
		return resp.containers, resp.err
	default:
		return nil, nil
	}
}

// ContainerInspect is the mock version of the actual method.
func (c *Client) ContainerInspect(ctx context.Context, container string) (types.ContainerJSON, error) {
	c.requests.containerInspect <- ContainerInspectRequest{container}
	select {
	case resp := <-c.responses.containerInspect:
		return resp.json, resp.err
	default:
		return types.ContainerJSON{}, nil
	}
}

// ContainerStop is the mock version of the actual method.
func (c *Client) ContainerStop(ctx context.Context, container string, timeout *time.Duration) error {
	c.requests.containerStop <- ContainerStopRequest{container}
	select {
	case resp := <-c.responses.containerStop:
		return resp.err
	default:
		return nil
	}
}

// ContainerRemove is the mock version of the actual method.
func (c *Client) ContainerRemove(ctx context.Context, container string, options types.ContainerRemoveOptions) error {
	c.requests.containerRemove <- ContainerRemoveRequest{container, options}
	select {
	case resp := <-c.responses.containerRemove:
		return resp.err
	default:
		return nil
	}
}

// ImageBuild is the mock version of the actual method.
func (c *Client) ImageBuild(ctx context.Context, context io.Reader,
	options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	bytes, err := ioutil.ReadAll(context)
	if err != nil {
		return types.ImageBuildResponse{}, err
	}
	c.requests.imageBuild <- ImageBuildRequest{bytes, options}
	select {
	case resp := <-c.responses.imageBuild:
		return resp.response, resp.err
	default:
		return types.ImageBuildResponse{}, nil
	}
}

// NetworkRemove is the mock version of the actual method.
func (c *Client) NetworkRemove(ctx context.Context, network string) error {
	c.requests.networkRemove <- NetworkRemoveRequest{network}
	return nil
}

// TaskList is the mock version of the actual method.
func (c *Client) TaskList(ctx context.Context, options types.TaskListOptions) ([]swarm.Task, error) {
	c.requests.taskList <- TaskListRequest{options}
	select {
	case resp := <-c.responses.taskList:
		return resp.tasks, resp.err
	default:
		return nil, nil
	}
}

// ServiceCreate is the mock version of the actual method.
func (c *Client) ServiceCreate(ctx context.Context,
	service swarm.ServiceSpec, options types.ServiceCreateOptions) (types.ServiceCreateResponse, error) {
	c.requests.serviceCreate <- ServiceCreateRequest{service, options}
	select {
	case resp := <-c.responses.serviceCreate:
		return resp.response, resp.err
	default:
		return types.ServiceCreateResponse{}, nil
	}
}

// ServiceList is the mock version of the actual method.
func (c *Client) ServiceList(ctx context.Context, options types.ServiceListOptions) ([]swarm.Service, error) {
	c.requests.serviceList <- ServiceListRequest{options}
	select {
	case resp := <-c.responses.serviceList:
		return resp.services, resp.err
	default:
		return nil, nil
	}
}

// ServiceInspectWithRaw is the mock version of the actual method.
func (c *Client) ServiceInspectWithRaw(ctx context.Context, serviceID string,
	options types.ServiceInspectOptions) (swarm.Service, []byte, error) {
	c.requests.serviceInspectWithRaw <- ServiceInspectWithRawRequest{serviceID, options}
	select {
	case resp := <-c.responses.serviceInspectWithRaw:
		return resp.service, resp.data, resp.err
	default:
		return swarm.Service{}, nil, nil
	}
}

// ServiceRemove is the mock version of the actual method.
func (c *Client) ServiceRemove(ctx context.Context, serviceID string) error {
	c.requests.serviceRemove <- ServiceRemoveRequest{serviceID}
	select {
	case resp := <-c.responses.serviceRemove:
		return resp.err
	default:
		return nil
	}
}

// ServiceLogs is the mock version of the actual method.
func (c *Client) ServiceLogs(ctx context.Context,
	serviceID string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	c.requests.serviceLogs <- ServiceLogsRequest{serviceID, options}
	select {
	case resp := <-c.responses.serviceLogs:
		return resp.rc, resp.err
	default:
		return nil, nil
	}
}

// Events is the mock version of the actual method.
func (c *Client) Events(ctx context.Context, options types.EventsOptions) (<-chan events.Message, <-chan error) {
	c.requests.events <- EventsRequest{Options: options}
	resp := <-c.responses.events
	return resp.messageChan, resp.errChan
}
