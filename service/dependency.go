package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/swarm"
	docker "github.com/fsouza/go-dockerclient"
)

func extractPorts(dependency *Dependency) (ports []swarm.PortConfig) {
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

// Start will start a dependency container
func (dependency *Dependency) Start(serviceName string, namespace string) (err error) {
	dependency.SwarmService, err = dockerCli.CreateService(docker.CreateServiceOptions{
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
	return
}

// Stop a dependency
func (dependency *Dependency) Stop(serviceName string, namespace string) (err error) {
	ctx := context.Background()
	dockerServices, err := dockerCli.ListServices(docker.ListServicesOptions{
		Filters: map[string][]string{
			"name": []string{strings.Join([]string{namespace, serviceName}, "_")},
		},
		Context: ctx,
	})
	if err != nil {
		return
	}
	dockerService := dockerServices[0]
	err = dockerCli.RemoveService(docker.RemoveServiceOptions{
		ID:      dockerService.ID,
		Context: ctx,
	})
	return
}
