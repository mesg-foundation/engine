package container

import (
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/docker/docker/client"
)

// ListServices returns existing docker services matching a specific label name.
func (c *DockerContainer) ListServices(labels ...string) ([]swarm.Service, error) {
	args := make([]filters.KeyValuePair, 0)
	for _, label := range labels {
		args = append(args, filters.KeyValuePair{
			Key:   "label",
			Value: label,
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	return c.client.ServiceList(ctx, types.ServiceListOptions{
		Filters: filters.NewArgs(args...),
	})
}

// FindService returns the Docker Service or an error if not found.
func (c *DockerContainer) FindService(namespace []string) (swarm.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	service, _, err := c.client.ServiceInspectWithRaw(ctx, c.Namespace(namespace),
		types.ServiceInspectOptions{},
	)
	return service, err
}

// StartService starts a docker service.
func (c *DockerContainer) StartService(options ServiceOptions) (serviceID string, err error) {
	service := options.toSwarmServiceSpec(c)
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	response, err := c.client.ServiceCreate(ctx, service, types.ServiceCreateOptions{})
	if err != nil {
		return "", err
	}
	return response.ID, c.waitForStatus(options.Namespace, RUNNING)
}

// StopService stops a docker service.
func (c *DockerContainer) StopService(namespace []string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	if err := c.client.ServiceRemove(ctx, c.Namespace(namespace)); err != nil && !docker.IsErrNotFound(err) {
		return err
	}
	if err := c.deletePendingContainer(namespace); err != nil {
		return err
	}
	return c.waitForStatus(namespace, STOPPED)
}

func (c *DockerContainer) deletePendingContainer(namespace []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	container, err := c.FindContainer(namespace)
	if err != nil {
		if docker.IsErrNotFound(err) {
			return nil
		}
		return err
	}
	// TOFIX: Hack to force Docker to remove the containers.
	// Sometime, the ServiceRemove function doesn't remove the associated containers.
	// This hack for Docker to stop and then remove the container.
	// See issue https://github.com/moby/moby/issues/32620
	if container.ContainerJSONBase != nil {
		timeout := 1 * time.Second
		c.client.ContainerStop(ctx, container.ID, &timeout)
		c.client.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{})
	}
	return c.deletePendingContainer(namespace)
}

// ServiceLogs returns the logs of a service.
func (c *DockerContainer) ServiceLogs(namespace []string) (io.ReadCloser, error) {
	return c.client.ServiceLogs(context.Background(), c.Namespace(namespace),
		types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Timestamps: false,
			Follow:     true,
		},
	)
}
