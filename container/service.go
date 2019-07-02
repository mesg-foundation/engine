package container

import (
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	return c.client.ServiceList(ctx, types.ServiceListOptions{
		Filters: filters.NewArgs(args...),
	})
}

// FindService returns the Docker Service or an error if not found.
func (c *DockerContainer) FindService(namespace string) (swarm.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	service, _, err := c.client.ServiceInspectWithRaw(ctx, c.Namespace(namespace),
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
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	response, err := c.client.ServiceCreate(ctx, service, types.ServiceCreateOptions{})
	if err != nil {
		return "", err
	}
	return response.ID, c.waitForStatus(options.Namespace, RUNNING)
}

// StopService stops a docker service.
func (c *DockerContainer) StopService(namespace string) error {
	status, err := c.Status(namespace)
	if err != nil {
		return err
	}
	if status == STOPPED {
		return nil
	}
	service, err := c.FindService(namespace)
	if err != nil && !docker.IsErrNotFound(err) {
		return err
	}
	stopGracePeriod := c.defaultStopGracePeriod
	if service.Spec.TaskTemplate.ContainerSpec != nil && service.Spec.TaskTemplate.ContainerSpec.StopGracePeriod != nil {
		stopGracePeriod = *service.Spec.TaskTemplate.ContainerSpec.StopGracePeriod
	}
	ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
	defer cancel()
	if err := c.client.ServiceRemove(ctx, c.Namespace(namespace)); err != nil && !docker.IsErrNotFound(err) {
		return err
	}
	if err := c.deletePendingContainer(namespace, time.Now().Add(stopGracePeriod)); err != nil {
		return err
	}
	return c.waitForStatus(namespace, STOPPED)
}

func (c *DockerContainer) deletePendingContainer(namespace string, maxGraceTime time.Time) error {
	container, err := c.FindContainer(namespace)
	if docker.IsErrNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	// Hack to force Docker to remove the containers.
	// Sometime, the ServiceRemove function doesn't remove the associated containers, or too late and the associated networks cannot be removed.
	// This hack for Docker to stop and then remove the container.
	// See issue https://github.com/moby/moby/issues/32620
	if time.Now().After(maxGraceTime) {
		ctx, cancel := context.WithTimeout(context.Background(), c.callTimeout)
		defer cancel()
		c.client.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{Force: true})
	}
	time.Sleep(1 * time.Second)
	return c.deletePendingContainer(namespace, maxGraceTime)
}
