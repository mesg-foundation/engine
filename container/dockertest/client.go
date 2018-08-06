package dockertest

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/docker/docker/api/types"
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

// requests holds request channels that holds call arguments of Docker client methods.
// Each call to a Docker client method piped to its request chan so multiple
// call to a same method can be received by reading its chan multiple times.
// Inspecting call arguments can be made by listening chan returned by LastX methods of *Testing.
type requests struct {
	negotiateAPIVersion   chan NegotiateAPIVersionRequest
	swarmInit             chan SwarmInitRequest
	networkCreate         chan NetworkCreateRequest
	info                  chan InfoRequest
	containerList         chan ContainerListRequest
	containerInspect      chan ContainerInspectRequest
	imageBuild            chan ImageBuildRequest
	networkInspect        chan NetworkInspectRequest
	networkRemove         chan NetworkRemoveRequest
	taskList              chan TaskListRequest
	serviceCreate         chan ServiceCreateRequest
	serviceList           chan ServiceListRequest
	serviceInspectWithRaw chan ServiceInspectWithRawRequest
	serviceRemove         chan ServiceRemoveRequest
	serviceLogs           chan ServiceLogsRequest
}

// responses holds response channels that holds 'faked' return values of Docker client methods.
// To fake a Docker client method's response send fake return values to it's
// response channel. This is done via ProvideX methods of *Testing.
// We use channels here instead of setting values on the struct in case of a need for returning
// conditional responses depending on method paramaters in future to deal with parallel calls to
// same client methods.
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
	contanerInspect       chan containerInspectResponse
	containerList         chan containerListResponse
}

// newClient returns a new mock Client for Docker.
func newClient() *Client {
	return &Client{
		// buffered channels helps with writing tests in synchronous syntax while using
		// LastX methods of Testing.
		requests: &requests{
			negotiateAPIVersion:   make(chan NegotiateAPIVersionRequest, 1),
			swarmInit:             make(chan SwarmInitRequest, 1),
			networkCreate:         make(chan NetworkCreateRequest, 1),
			info:                  make(chan InfoRequest, 1),
			containerList:         make(chan ContainerListRequest, 1),
			containerInspect:      make(chan ContainerInspectRequest, 1),
			imageBuild:            make(chan ImageBuildRequest, 1),
			networkInspect:        make(chan NetworkInspectRequest, 1),
			networkRemove:         make(chan NetworkRemoveRequest, 1),
			taskList:              make(chan TaskListRequest, 1),
			serviceCreate:         make(chan ServiceCreateRequest, 1),
			serviceList:           make(chan ServiceListRequest, 1),
			serviceInspectWithRaw: make(chan ServiceInspectWithRawRequest, 1),
			serviceRemove:         make(chan ServiceRemoveRequest, 1),
			serviceLogs:           make(chan ServiceLogsRequest, 1),
		},

		// buffered channels helps with writing tests in synchronous syntax while using
		// ProvideX methods of Testing.
		responses: &responses{
			info:                  make(chan infoResponse, 1),
			imageBuild:            make(chan imageBuildResponse, 1),
			networkInspect:        make(chan networkInspectResponse, 1),
			networkCreate:         make(chan networkCreateResponse, 1),
			taskList:              make(chan taskListResponse, 1),
			serviceCreate:         make(chan serviceCreateResponse, 1),
			serviceList:           make(chan serviceListResponse, 1),
			serviceInspectWithRaw: make(chan serviceInspectWithRawResponse, 1),
			serviceRemove:         make(chan serviceRemoveResponse, 1),
			serviceLogs:           make(chan serviceLogsResponse, 1),
			contanerInspect:       make(chan containerInspectResponse, 1),
			containerList:         make(chan containerListResponse, 1),
		},
	}
}

// note that select: default statements are a short cut of returning zero value responses
// instead of a explicit need to call Provide methods with zero values.

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
	case resp := <-c.responses.contanerInspect:
		return resp.json, resp.err
	default:
		return types.ContainerJSON{}, nil
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
