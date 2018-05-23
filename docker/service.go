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
func FindService(name []string) (dockerService *swarm.Service, err error) {
	ctx := context.Background()
	client, err := Client()
	if err != nil {
		return
	}
	dockerServices, err := client.ListServices(godocker.ListServicesOptions{
		Filters: map[string][]string{
			"name": []string{Namespace(name)},
		},
		Context: ctx,
	})
	if err != nil {
		return
	}
	dockerService = serviceMatch(dockerServices, name)
	return
}

func serviceMatch(dockerServices []swarm.Service, name []string) (dockerService *swarm.Service) {
	for _, service := range dockerServices {
		if service.Spec.Annotations.Name == Namespace(name) {
			dockerService = &service
			break
		}
	}
	return
}

// StartService starts a docker service
func StartService(service godocker.CreateServiceOptions) (dockerService *swarm.Service, err error) {
	client, err := Client()
	if err != nil {
		return
	}
	return client.CreateService(service)
}

// StopService stops a docker service
func StopService(name []string) (err error) {
	client, err := Client()
	if err != nil {
		return
	}
	if !IsServiceRunning(name) {
		return
	}
	dockerService, err := FindService(name)
	if err == nil && dockerService.ID != "" {
		err = client.RemoveService(godocker.RemoveServiceOptions{
			ID:      dockerService.ID,
			Context: context.Background(),
		})
	}
	return
}

// ServiceStatus return the status of the Docker Swarm Servicer
func ServiceStatus(name []string) (status StatusType) {
	dockerService, err := FindService(name)
	status = STOPPED
	if err == nil && dockerService != nil && dockerService.ID != "" {
		status = RUNNING
	}
	return
}

// IsServiceRunning returns true if the dependency is running, false otherwise
func IsServiceRunning(name []string) (running bool) {
	running = ServiceStatus(name) == RUNNING
	return
}

// IsServiceStopped returns true if the dependency is stopped, false otherwise
func IsServiceStopped(name []string) (running bool) {
	running = ServiceStatus(name) == STOPPED
	return
}
