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

	lastNegotiateAPIVersion chan struct{}
	lastSwarmInit           chan swarm.InitRequest
	lastNetworkCreate       chan NetworkCreate
	lastInfo                chan types.Info
	lastContainerList       chan types.ContainerListOptions
	lastContainerInspect    chan string
	lastImageBuild          chan ImageBuild
	lastNetworkInspect      chan NetworkInspect
	lastNetworkRemove       chan string

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
}

func newClient() *client {
	return &client{
		lastNegotiateAPIVersion: make(chan struct{}, 1),
		lastSwarmInit:           make(chan swarm.InitRequest, 1),
		lastNetworkCreate:       make(chan NetworkCreate, 1),
		lastInfo:                make(chan types.Info, 1),
		lastContainerList:       make(chan types.ContainerListOptions, 1),
		lastContainerInspect:    make(chan string, 1),
		lastImageBuild:          make(chan ImageBuild, 1),
		lastNetworkInspect:      make(chan NetworkInspect, 1),
		lastNetworkRemove:       make(chan string, 1),
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
	return c.containers, ErrContainerDoesNotExists
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
