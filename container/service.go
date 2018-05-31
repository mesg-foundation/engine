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
func StartService(spec swarm.ServiceSpec) (serviceID string, err error) {
	client, err := Client()
	if err != nil {
		return
	}
	response, err := client.ServiceCreate(context.Background(), spec, types.ServiceCreateOptions{})
	if err != nil {
		return
	}
	serviceID = response.ID
	return
}

// StopService stops a docker service
func StopService(namespace []string) (err error) {
	stopped, err := IsServiceStopped(namespace)
	if err != nil || stopped == true {
		return
	}
	client, err := Client()
	if err != nil {
		return
	}
	err = client.ServiceRemove(context.Background(), Namespace(namespace))
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

// IsServiceRunning returns true if the service is running, false otherwise
func IsServiceRunning(namespace []string) (result bool, err error) {
	status, err := ServiceStatus(namespace)
	result = status == RUNNING
	return
}

// IsServiceStopped returns true if the service is stopped, false otherwise
func IsServiceStopped(namespace []string) (result bool, err error) {
	status, err := ServiceStatus(namespace)
	result = status == STOPPED
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
