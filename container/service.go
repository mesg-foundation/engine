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
func ListServices(label string) (services []swarm.Service, err error) {
	client, err := Client()
	if err != nil {
		return
	}
	services, err = client.ServiceList(context.Background(), types.ServiceListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: label,
		}),
	})
	return
}

// FindService returns the Docker Service. Return error if not found.
func FindService(namespace []string) (service swarm.Service, err error) {
	client, err := Client()
	if err != nil {
		return
	}
	service, _, err = client.ServiceInspectWithRaw(context.Background(), Namespace(namespace), types.ServiceInspectOptions{})
	return
}

// StartService starts a docker service
func StartService(options ServiceOptions) (serviceID string, err error) {
	client, err := Client()
	if err != nil {
		return
	}
	service := options.toSwarmServiceSpec()
	response, err := client.ServiceCreate(context.Background(), service, types.ServiceCreateOptions{})
	if err != nil {
		return
	}
	serviceID = response.ID
	err = waitForStatus(options.Namespace, RUNNING)
	return
}

// StopService stops a docker service
func StopService(namespace []string) (err error) {
	status, err := ServiceStatus(namespace)
	if err != nil || status == STOPPED {
		return
	}
	client, err := Client()
	if err != nil {
		return
	}
	err = client.ServiceRemove(context.Background(), Namespace(namespace))
	if err != nil {
		return
	}
	err = waitForStatus(namespace, STOPPED)
	return
}

// ServiceStatus return the status of the Docker Swarm Servicer
func ServiceStatus(namespace []string) (status StatusType, err error) {
	status = STOPPED
	_, err = FindService(namespace)
	if docker.IsErrNotFound(err) {
		err = nil
		return
	}
	status = RUNNING
	return
}

// ServiceLogs returns the logs of a service
func ServiceLogs(namespace []string) (reader io.ReadCloser, err error) {
	client, err := Client()
	if err != nil {
		return
	}
	reader, err = client.ServiceLogs(context.Background(), Namespace(namespace), types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: false,
		Follow:     true,
	})
	return
}
