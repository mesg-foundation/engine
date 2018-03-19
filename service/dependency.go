package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/swarm"
	docker "github.com/fsouza/go-dockerclient"
)

func extractPorts(dependency Dependency) (ports []swarm.PortConfig) {
	ports = make([]swarm.PortConfig, len(dependency.Ports))
	for i, p := range dependency.Ports {
		split := strings.Split(p, ":")
		from, _ := strconv.ParseUint(split[0], 10, 64)
		to := from
		if len(split) > 1 {
			to, _ = strconv.ParseUint(split[1], 10, 64)
		}
		ports[i] = swarm.PortConfig{
			Protocol:      swarm.PortConfigProtocolTCP,
			PublishMode:   swarm.PortConfigPublishModeIngress,
			TargetPort:    uint32(to),
			PublishedPort: uint32(from),
		}
	}
	return
}

func (dependency Dependency) getDockerService(namespace string, dependencyName string) (dockerService swarm.Service, err error) {
	ctx := context.Background()
	dockerServices, err := dockerCli.ListServices(docker.ListServicesOptions{
		Filters: map[string][]string{
			"name": []string{strings.Join([]string{namespace, dependencyName}, "_")},
		},
		Context: ctx,
	})
	if err != nil {
		return
	}
	dockerService = dockerServiceMatch(dockerServices, namespace, dependencyName)
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

// Stop a dependency
func (dependency Dependency) Stop(namespace string, dependencyName string) (err error) {
	ctx := context.Background()
	if !dependency.IsRunning(namespace, dependencyName) {
		return
	}
	dockerService, err := dependency.getDockerService(namespace, dependencyName)
	if err == nil && dockerService.ID != "" {
		err = dockerCli.RemoveService(docker.RemoveServiceOptions{
			ID:      dockerService.ID,
			Context: ctx,
		})
	}
	return
}
