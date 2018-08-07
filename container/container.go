package container

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/docker/docker/client"
)

// Container provides functionaliets for Docker containers for MESG.
type Container struct {
	// client is a Docker client.
	client docker.CommonAPIClient

	// callTimeout is the timeout value for Docker API calls.
	callTimeout time.Duration
}

// Option is a configuration func for Container.
type Option func(*Container)

// New creates a new Container with given options.
func New(options ...Option) (*Container, error) {
	c := &Container{
		callTimeout: time.Second * 10,
	}
	for _, option := range options {
		option(c)
	}
	var err error
	if c.client == nil {
		c.client, err = docker.NewClientWithOpts(docker.FromEnv)
		if err != nil {
			return c, err
		}
	}
	c.negotiateAPIVersion()
	if err := c.createSwarmIfNeeded(); err != nil {
		return c, err
	}
	return c, c.createSharedNetworkIfNeeded()
}

// ClientOption creates a new Option with given docker client for Container.
func ClientOption(client docker.CommonAPIClient) Option {
	return func(c *Container) {
		c.client = client
	}
}

// TimeoutOption creates a new Option with given d http call timeout for Container.
func TimeoutOption(d time.Duration) Option {
	return func(c *Container) {
		c.callTimeout = d
	}
}

func (c *Container) negotiateAPIVersion() {
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	c.client.NegotiateAPIVersion(ctx)
}

func (c *Container) createSwarmIfNeeded() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	info, err := c.client.Info(ctx)
	if err != nil {
		return err
	}
	if info.Swarm.NodeID != "" {
		return nil
	}
	ctx, cancel = context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	_, err = c.client.SwarmInit(ctx, swarm.InitRequest{
		ListenAddr: "0.0.0.0:2377", // https://docs.docker.com/engine/reference/commandline/swarm_init/#usage
	})
	return err
}

// FindContainer returns a docker container.
func (c *Container) FindContainer(namespace []string) (types.ContainerJSON, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	containers, err := c.client.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "com.docker.stack.namespace=" + Namespace(namespace),
		}),
		Limit: 1,
	})
	if err != nil {
		return types.ContainerJSON{}, err
	}
	containerID := ""
	if len(containers) == 1 {
		containerID = containers[0].ID
	}
	ctx, cancel = context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	return c.client.ContainerInspect(ctx, containerID)
}

// Status returns the status of a docker container.
func (c *Container) Status(namespace []string) (StatusType, error) {
	status := STOPPED
	container, err := c.FindContainer(namespace)
	if docker.IsErrNotFound(err) {
		return status, nil
	}
	if err != nil {
		return status, err
	}
	if container.State.Running {
		status = RUNNING
	}
	return status, nil
}
