package docker

import (
	"context"

	"github.com/docker/docker/api/types/swarm"
	godocker "github.com/fsouza/go-dockerclient"
)

// StatusType of the service
type StatusType uint

// status for services
const (
	STOPPED StatusType = 0
	RUNNING StatusType = 1
	PARTIAL StatusType = 2
)

//  FindService returns the Docker Service
func FindService(namespace []string) (dockerService *swarm.Service, err error) {
	ctx := context.Background()
	client, err := Client()
	if err != nil {
		return
	}
	dockerServices, err := client.ListServices(godocker.ListServicesOptions{
		Filters: map[string][]string{
			"name": []string{Namespace(namespace)},
		},
		Context: ctx,
	})
	if err != nil {
		return
	}
	dockerService = serviceMatch(dockerServices, namespace)
	return
}

func serviceMatch(dockerServices []swarm.Service, namespace []string) *swarm.Service {
	for _, service := range dockerServices {
		if service.Spec.Annotations.Name == Namespace(namespace) {
			return &service
		}
	}
	return nil
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
	if !IsServiceRunning(namespace) {
		return
	}
	dockerService, err := FindService(namespace)
	if err == nil && dockerService.ID != "" {
		err = client.RemoveService(godocker.RemoveServiceOptions{
			ID:      dockerService.ID,
			Context: context.Background(),
		})
	}
	return
}

// ServiceStatus return the status of the Docker Swarm Servicer
func ServiceStatus(namespace []string) (status StatusType) {
	dockerService, err := FindService(namespace)
	status = STOPPED
	if err == nil && dockerService != nil && dockerService.ID != "" {
		status = RUNNING
	}
	return
}

// IsServiceRunning returns true if the dependency is running, false otherwise
func IsServiceRunning(namespace []string) (running bool) {
	running = ServiceStatus(namespace) == RUNNING
	return
}

// IsServiceStopped returns true if the dependency is stopped, false otherwise
func IsServiceStopped(namespace []string) (running bool) {
	running = ServiceStatus(namespace) == STOPPED
	return
}
