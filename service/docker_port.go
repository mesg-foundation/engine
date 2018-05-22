package service

import (
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/swarm"
)

// DockerPorts extract ports from a Dependency and transform them to a swarm.PortConfig
func DockerPorts(dependency *Dependency) (ports []swarm.PortConfig) {
	ports = make([]swarm.PortConfig, len(dependency.GetPorts()))
	for i, p := range dependency.GetPorts() {
		split := strings.Split(p, ":")
		published, _ := strconv.ParseUint(split[0], 10, 64)
		target := published
		if len(split) > 1 {
			target, _ = strconv.ParseUint(split[1], 10, 64)
		}
		ports[i] = swarm.PortConfig{
			Protocol:      swarm.PortConfigProtocolTCP,
			PublishMode:   swarm.PortConfigPublishModeIngress,
			TargetPort:    uint32(target),
			PublishedPort: uint32(published),
		}
	}
	return
}
