package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/swarm"
	docker "github.com/fsouza/go-dockerclient"
)

func extractPorts(dependency *Dependency) (ports []swarm.PortConfig) {
	ports = make([]swarm.PortConfig, len(dependency.Ports))
	for i, p := range dependency.Ports {
		split := strings.Split(p, ":")
		fmt.Println(split)
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
