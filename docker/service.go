package docker

import (
	"context"

	"github.com/docker/docker/api/types/swarm"
	godocker "github.com/fsouza/go-dockerclient"
)

// ListServices returns existing docker services matching a specific label name
func ListServices(label string) (dockerServices []swarm.Service, err error) {
	client, err := Client()
	if err != nil {
		return
	}
	dockerServices, err = client.ListServices(godocker.ListServicesOptions{
		Context: context.Background(),
		Filters: map[string][]string{
			"label": []string{label},
		},
	})
	if err != nil {
		return
	}
	return
}

//  FindService returns the Docker Service
func FindService(namespace []string) (dockerService *swarm.Service, err error) {
	client, err := Client()
	if err != nil {
		return
	}
	dockerServices, err := client.ListServices(godocker.ListServicesOptions{
		Filters: map[string][]string{
			"name": []string{Namespace(namespace)},
		},
		Context: context.Background(),
	})
	if err != nil {
		return
	}
	for _, service := range dockerServices {
		if service.Spec.Name == Namespace(namespace) {
			dockerService = &service
			return
		}
	}
	return
}

// StartService starts a docker service
func StartService(options *ServiceOptions) (dockerService *swarm.Service, err error) {
	client, err := Client()
	if err != nil {
		return
	}
	options.merge()
	return client.CreateService(*options.CreateServiceOptions)
}

// StopService stops a docker service
func StopService(namespace []string) (err error) {
	client, err := Client()
	if err != nil {
		return
	}
	stopped, err := IsServiceStopped(namespace)
	if err != nil || stopped == true {
		return
	}
	dockerService, err := FindService(namespace)
	if err == nil && dockerService != nil && dockerService.ID != "" {
		err = client.RemoveService(godocker.RemoveServiceOptions{
			ID:      dockerService.ID,
			Context: context.Background(),
		})
	}
	return
}

// ServiceStatus return the status of the Docker Swarm Servicer
func ServiceStatus(namespace []string) (status StatusType, err error) {
	dockerService, err := FindService(namespace)
	if err != nil {
		return
	}
	status = STOPPED
	if err == nil && dockerService != nil && dockerService.ID != "" {
		status = RUNNING
	}
	return
}

// IsServiceRunning returns true if the service is running, false otherwise
func IsServiceRunning(namespace []string) (result bool, err error) {
	status, err := ServiceStatus(namespace)
	if err != nil {
		return
	}
	result = status == RUNNING
	return
}

// IsServiceStopped returns true if the service is stopped, false otherwise
func IsServiceStopped(namespace []string) (result bool, err error) {
	status, err := ServiceStatus(namespace)
	if err != nil {
		return
	}
	result = status == STOPPED
	return
}
