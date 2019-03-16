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
	args := make([]filters.KeyValuePair, len(labels))
	for i, label := range labels {
		args[i] = filters.KeyValuePair{
			Key:   "label",
			Value: label,
		}
	}
	return c.client.ServiceList(context.Background(), types.ServiceListOptions{
		Filters: filters.NewArgs(args...),
	})
}

// FindService returns the Docker Service or an error if not found.
func (c *DockerContainer) FindService(namespace []string) (swarm.Service, error) {
	service, _, err := c.client.ServiceInspectWithRaw(context.Background(), c.Namespace(namespace),
		types.ServiceInspectOptions{},
	)
	return service, err
}

// StartService starts a docker service.
func (c *DockerContainer) StartService(options ServiceOptions) (serviceID string, err error) {
	status, err := c.Status(options.Namespace)
	if err != nil {
		return "", err
	}
	if status == RUNNING {
		service, err := c.FindService(options.Namespace)
		return service.ID, err
	}

	service := options.toSwarmServiceSpec(c)
	response, err := c.client.ServiceCreate(context.Background(), service, types.ServiceCreateOptions{})
	if err != nil {
		return "", err
	}
	return response.ID, c.waitForStatus(options.Namespace, RUNNING)
}

// StopService stops a docker service.
func (c *DockerContainer) StopService(namespace []string) error {
	status, err := c.Status(namespace)
	if err != nil {
		return err
	}
	if status == STOPPED {
		return nil
	}

	if err := c.client.ServiceRemove(context.Background(), c.Namespace(namespace)); err != nil && !docker.IsErrNotFound(err) {
		return err
	}
	if err := c.deletePendingContainer(namespace); err != nil {
		return err
	}
	return c.waitForStatus(namespace, STOPPED)
}

func (c *DockerContainer) deletePendingContainer(namespace []string) error {
	container, err := c.FindContainer(namespace)
	if docker.IsErrNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	// TOFIX: Hack to force Docker to remove the containers.
	// Sometime, the ServiceRemove function doesn't remove the associated containers.
	// This hack for Docker to stop and then remove the container.
	// See issue https://github.com/moby/moby/issues/32620
	c.client.ContainerStop(context.Background(), container.ID, nil)
	c.client.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{})
	time.Sleep(1 * time.Second)
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
