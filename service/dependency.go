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

func getDockerService(namespace string, dependencyName string) (dockerService swarm.Service, err error) {
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

func dependencyStatus(namespace string, dependencyName string) (status StatusType) {
	ctx := context.Background()
	dockerServices, err := dockerCli.ListServices(docker.ListServicesOptions{
		Filters: map[string][]string{
			"name": []string{strings.Join([]string{namespace, dependencyName}, "_")},
		},
		Context: ctx,
	})
	dockerService := dockerServiceMatch(dockerServices, namespace, dependencyName)
	status = STOPPED
	if err == nil && dockerService.ID != "" {
		status = RUNNING
	}
	return
}
