package dockertest

import (
	"context"
	"io"
	"io/ioutil"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/docker/docker/client"
)

type client struct {
	// a placeholder for unimplemented methods.
	docker.CommonAPIClient

	lastNegotiateAPIVersion   chan struct{}
	lastSwarmInit             chan swarm.InitRequest
	lastNetworkCreate         chan NetworkCreate
	lastInfo                  chan types.Info
	lastContainerList         chan types.ContainerListOptions
	lastContainerInspect      chan string
	lastImageBuild            chan ImageBuild
	lastNetworkInspect        chan NetworkInspect
	lastNetworkRemove         chan string
	lastTaskList              chan types.TaskListOptions
	lastServiceCreate         chan ServiceCreate
	lastServiceList           chan types.ServiceListOptions
	lastServiceInspectWithRaw chan ServiceInspectWithRaw
	lastServiceRemove         chan string
	lastServiceLogs           chan ServiceLogs

	m sync.RWMutex

	info    types.Info
	infoErr error

	containers       []types.Container
	containerInspect types.ContainerJSON

	imageBuild types.ImageBuildResponse

	networkInspect    types.NetworkResource
	networkInspectErr error

	networkCreate    types.NetworkCreateResponse
	networkCreateErr error

	tasklist    []swarm.Task
	tasklistErr error

	serviceCreate    types.ServiceCreateResponse
	serviceCreateErr error

	serviceList    []swarm.Service
	serviceListErr error

	serviceInspectWithRaw      swarm.Service
	serviceInspectWithRawBytes []byte
	serviceInspectWithRawErr   error

	serviceRemoveErr error

	serviceLogs    io.ReadCloser
	serviceLogsErr error
}

func newClient() *client {
	return &client{
		lastNegotiateAPIVersion:   make(chan struct{}, 1),
		lastSwarmInit:             make(chan swarm.InitRequest, 1),
		lastNetworkCreate:         make(chan NetworkCreate, 1),
		lastInfo:                  make(chan types.Info, 1),
		lastContainerList:         make(chan types.ContainerListOptions, 1),
		lastContainerInspect:      make(chan string, 1),
		lastImageBuild:            make(chan ImageBuild, 1),
		lastNetworkInspect:        make(chan NetworkInspect, 1),
		lastNetworkRemove:         make(chan string, 1),
		lastTaskList:              make(chan types.TaskListOptions, 1),
		lastServiceCreate:         make(chan ServiceCreate, 1),
		lastServiceList:           make(chan types.ServiceListOptions, 1),
		lastServiceInspectWithRaw: make(chan ServiceInspectWithRaw, 1),
		lastServiceRemove:         make(chan string, 1),
		lastServiceLogs:           make(chan ServiceLogs, 1),
	}
}

func (c *client) NegotiateAPIVersion(ctx context.Context) {
	c.lastNegotiateAPIVersion <- struct{}{}
}

type NetworkInspect struct {
	Network string
	Options types.NetworkInspectOptions
}

func (c *client) NetworkInspect(ctx context.Context, network string,
	options types.NetworkInspectOptions) (types.NetworkResource, error) {
	c.lastNetworkInspect <- NetworkInspect{network, options}
	return c.networkInspect, c.networkInspectErr
}

type NetworkCreate struct {
	Name    string
	Options types.NetworkCreate
}

func (c *client) NetworkCreate(ctx context.Context, name string,
	options types.NetworkCreate) (types.NetworkCreateResponse, error) {
	c.lastNetworkCreate <- NetworkCreate{name, options}
	return c.networkCreate, c.networkCreateErr
}

func (c *client) SwarmInit(ctx context.Context, req swarm.InitRequest) (string, error) {
	c.lastSwarmInit <- req
	return "", nil
}

func (c *client) Info(context.Context) (types.Info, error) {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.info, c.infoErr
}

func (c *client) ContainerList(ctx context.Context,
	options types.ContainerListOptions) ([]types.Container, error) {
	c.lastContainerList <- options
	if len(c.containers) > 0 {
		return c.containers, nil
	}
	return c.containers, NotFoundErr{}
}

func (c *client) ContainerInspect(ctx context.Context, container string) (types.ContainerJSON, error) {
	c.lastContainerInspect <- container
	return c.containerInspect, nil
}

type ImageBuild struct {
	FileData []byte
	Options  types.ImageBuildOptions
}

func (c *client) ImageBuild(ctx context.Context, context io.Reader, options types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	bytes, err := ioutil.ReadAll(context)
	if err != nil {
		return types.ImageBuildResponse{}, err
	}
	c.lastImageBuild <- ImageBuild{bytes, options}
	return c.imageBuild, nil
}

func (c *client) NetworkRemove(ctx context.Context, network string) error {
	c.lastNetworkRemove <- network
	return nil
}

func (c *client) TaskList(ctx context.Context, options types.TaskListOptions) ([]swarm.Task, error) {
	c.lastTaskList <- options
	return c.tasklist, c.tasklistErr
}

type ServiceCreate struct {
	Service swarm.ServiceSpec
	Options types.ServiceCreateOptions
}

func (c *client) ServiceCreate(ctx context.Context,
	service swarm.ServiceSpec, options types.ServiceCreateOptions) (types.ServiceCreateResponse, error) {
	c.lastServiceCreate <- ServiceCreate{service, options}
	return c.serviceCreate, c.serviceCreateErr
}

func (c *client) ServiceList(ctx context.Context, options types.ServiceListOptions) ([]swarm.Service, error) {
	c.lastServiceList <- options
	return c.serviceList, c.serviceListErr
}

type ServiceInspectWithRaw struct {
	ServiceID string
	Options   types.ServiceInspectOptions
}

func (c *client) ServiceInspectWithRaw(ctx context.Context, serviceID string,
	options types.ServiceInspectOptions) (swarm.Service, []byte, error) {
	c.lastServiceInspectWithRaw <- ServiceInspectWithRaw{serviceID, options}
	return c.serviceInspectWithRaw, c.serviceInspectWithRawBytes, c.serviceInspectWithRawErr
}

func (c *client) ServiceRemove(ctx context.Context, serviceID string) error {
	c.lastServiceRemove <- serviceID
	return c.serviceRemoveErr
}

type ServiceLogs struct {
	ServiceID string
	Options   types.ContainerLogsOptions
}

func (c *client) ServiceLogs(ctx context.Context,
	serviceID string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	c.lastServiceLogs <- ServiceLogs{serviceID, options}
	return c.serviceLogs, c.serviceLogsErr
}
