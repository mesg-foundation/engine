package container

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/docker/docker/client"
	"github.com/mesg-foundation/core/x/xstrings"
)

// Container provides high level interactions with Docker API for MESG.
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
		callTimeout: 10 * time.Second,
	}
	for _, option := range options {
		option(c)
	}
	var err error
	if c.client == nil {
		c.client, err = docker.NewEnvClient()
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

// ClientOption receives a client which will be used to interact with Docker API.
func ClientOption(client docker.CommonAPIClient) Option {
	return func(c *Container) {
		c.client = client
	}
}

// TimeoutOption receives d which will be set as a timeout value for Docker API calls.
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

// Status returns the status of the container based on the docker container and docker service
// If any error happen the status will be `UNKNOWN`
// Otherwise the following status will be applied:
//   - RUNNING: When the container is running in docker (whatever the status of the service)
//   - STARTING: When the service is running but the container is not yet started
//   - STOPPED: When neither the container nor the service are running in docker
func (c *Container) Status(namespace []string) (StatusType, error) {
	container, err := c.containerExists(namespace)
	if err != nil {
		return UNKNOWN, err
	}
	service, err := c.serviceExists(namespace)
	if err != nil {
		return UNKNOWN, err
	}

	statuses := []struct {
		container bool
		service   bool
		status    StatusType
	}{
		{service: true, container: true, status: RUNNING},
		{service: true, container: false, status: STARTING},
		{service: false, container: true, status: RUNNING}, // This is actually stopping
		{service: false, container: false, status: STOPPED},
	}

	for _, s := range statuses {
		if s.container == container && s.service == service {
			return s.status, nil
		}
	}
	return UNKNOWN, nil // This should never be reached but it's better than a panic :)
}

// containerExists return true if the docker service can be found, return false otherwise and return
// any errors that are not `NotFound` errors
func (c *Container) containerExists(namespace []string) (bool, error) {
	container, err := c.FindContainer(namespace)
	if err == nil && container.State != nil &&
		xstrings.SliceContains([]string{"exited", "dead"}, container.State.Status) {
		return false, nil
	}
	return presenceHandling(err)
}

// serviceExists return true if the docker container can be found, return false otherwise and return
// any errors that are not `NotFound` errors
func (c *Container) serviceExists(namespace []string) (bool, error) {
	_, err := c.FindService(namespace)
	return presenceHandling(err)
}

// presenceHandling handle the error the check the presence of a docker resource.
// It returns any error that are not docker `NotFound` errors
// It returns a boolean that says if the docker resource exists
func presenceHandling(err error) (bool, error) {
	if err != nil && !docker.IsErrNotFound(err) {
		return false, err
	}
	return !docker.IsErrNotFound(err), nil
}
