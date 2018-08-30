package container

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/docker/docker/client"
)

// ListServices returns existing docker services matching a specific label name.
func (c *Container) ListServices(label string) ([]swarm.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	return c.client.ServiceList(ctx, types.ServiceListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: label,
		}),
	})
}

// FindService returns the Docker Service or an error if not found.
func (c *Container) FindService(namespace []string) (swarm.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	service, _, err := c.client.ServiceInspectWithRaw(ctx, Namespace(namespace),
		types.ServiceInspectOptions{},
	)
	return service, err
}

// StartService starts a docker service.
func (c *Container) StartService(options ServiceOptions) (serviceID string, err error) {
	service := options.toSwarmServiceSpec()
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	response, err := c.client.ServiceCreate(ctx, service, types.ServiceCreateOptions{})
	if err != nil {
		return "", err
	}
	return response.ID, c.waitForStatus(options.Namespace, RUNNING)
}

// StopService stops a docker service.
func (c *Container) StopService(namespace []string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	if err := c.client.ServiceRemove(ctx, Namespace(namespace)); err != nil {
		if docker.IsErrNotFound(err) {
			return nil
		}
		return err
	}
	return c.waitForStatus(namespace, STOPPED)
}

// ServiceLogs returns the logs of a service.
func (c *Container) ServiceLogs(namespace []string) (io.ReadCloser, error) {
	return c.client.ServiceLogs(context.Background(), Namespace(namespace),
		types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Timestamps: false,
			Follow:     true,
		},
	)
}
