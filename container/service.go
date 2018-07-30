package container

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/docker/docker/client"
)

// ListServices returns existing docker services matching a specific label name
func ListServices(label string) ([]swarm.Service, error) {
	client, err := Client()
	if err != nil {
		return nil, err
	}
	return client.ServiceList(context.Background(), types.ServiceListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: label,
		}),
	})
}

// FindService returns the Docker Service. Return error if not found.
func FindService(namespace []string) (swarm.Service, error) {
	client, err := Client()
	if err != nil {
		return swarm.Service{}, err
	}
	service, _, err := client.ServiceInspectWithRaw(context.Background(), Namespace(namespace), types.ServiceInspectOptions{})
	return service, err
}

// StartService starts a docker service
func StartService(options ServiceOptions) (string, error) {
	client, err := Client()
	if err != nil {
		return "", err
	}
	service := options.toSwarmServiceSpec()
	response, err := client.ServiceCreate(context.Background(), service, types.ServiceCreateOptions{})
	if err != nil {
		return "", err
	}
	return response.ID, waitForStatus(options.Namespace, RUNNING)
}

// StopService stops a docker service
func StopService(namespace []string) error {
	status, err := ServiceStatus(namespace)
	if err != nil || status == STOPPED {
		return err
	}
	client, err := Client()
	if err != nil {
		return err
	}
	if err := client.ServiceRemove(context.Background(), Namespace(namespace)); err != nil {
		return err
	}
	return waitForStatus(namespace, STOPPED)
}

// ServiceStatus return the status of the Docker Swarm Servicer
func ServiceStatus(namespace []string) (StatusType, error) {
	if _, err := FindService(namespace); !docker.IsErrNotFound(err) {
		return STOPPED, err
	}
	return RUNNING, nil
}

// ServiceLogs returns the logs of a service
func ServiceLogs(namespace []string) (io.ReadCloser, error) {
	client, err := Client()
	if err != nil {
		return nil, err
	}
	return client.ServiceLogs(context.Background(), Namespace(namespace), types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: false,
		Follow:     true,
	})
}
