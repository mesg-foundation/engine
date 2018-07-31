package dockertest

import (
	"context"
	"errors"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/docker/docker/client"
)

var (
	ErrContainerDoesNotExists = errors.New("containers does not exists")
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

	m sync.RWMutex

	info    types.Info
	infoErr error

	containers       []types.Container
	containerInspect types.ContainerJSON
}

func newClient() *client {
	return &client{
		lastNegotiateAPIVersion: make(chan struct{}, 1),
		lastSwarmInit:           make(chan swarm.InitRequest, 1),
		lastNetworkCreate:       make(chan NetworkCreate, 1),
		lastInfo:                make(chan types.Info, 1),
		lastContainerList:       make(chan types.ContainerListOptions, 1),
		lastContainerInspect:    make(chan string, 1),
	}
}

func (c *client) NegotiateAPIVersion(ctx context.Context) {
	c.lastNegotiateAPIVersion <- struct{}{}
}

func (c *client) NetworkInspect(ctx context.Context, network string,
	options types.NetworkInspectOptions) (types.NetworkResource, error) {
	return types.NetworkResource{}, nil
}

type NetworkCreate struct {
	Name    string
	Options types.NetworkCreate
}

func (c *client) NetworkCreate(ctx context.Context, name string,
	options types.NetworkCreate) (types.NetworkCreateResponse, error) {
	c.lastNetworkCreate <- NetworkCreate{name, options}
	return types.NetworkCreateResponse{}, nil
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
