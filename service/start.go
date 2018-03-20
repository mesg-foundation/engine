package service

import (
	"strings"

	"github.com/docker/docker/api/types/swarm"
	docker "github.com/fsouza/go-dockerclient"
)

// Start a service
func (service *Service) Start() (dockerServices []*swarm.Service, err error) {
	if service.IsRunning() {
		return
	}
	// If there is one but not all services running stop to restart all
	if service.IsPartiallyRunning() {
		service.Stop()
	}
	dockerServices = make([]*swarm.Service, len(service.Dependencies))
	i := 0
	for name, dependency := range service.Dependencies {
		dockerServices[i], err = dependency.Start(service.namespace(), name)
		i++
		if err != nil {
			break
		}
	}
	// Disgrasfully close the service because there is an error
	if err != nil {
		service.Stop()
	}
	return
}

// Start will start a dependency container
func (dependency Dependency) Start(namespace string, serviceName string) (dockerService *swarm.Service, err error) {
	return dockerCli.CreateService(docker.CreateServiceOptions{
		ServiceSpec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: strings.Join([]string{namespace, serviceName}, "_"),
				Labels: map[string]string{
					"labelImage":     dependency.Image,
					"labelNamespace": namespace,
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: &swarm.ContainerSpec{
					Image: dependency.Image,
					Args:  strings.Fields(dependency.Command),
					Labels: map[string]string{
						"labelNamespace": namespace,
					},
				},
			},
			EndpointSpec: &swarm.EndpointSpec{
				Ports: extractPorts(dependency),
			},
		},
	})
}
