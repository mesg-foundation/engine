package service

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/mesg-foundation/core/config"
	"github.com/spf13/viper"
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

func extractVolumes(service *Service, dependency *Dependency, details dependencyDetails) (volumes []mount.Mount) {
	volumes = make([]mount.Mount, 0)
	for _, volume := range dependency.Volumes {
		path := filepath.Join(details.namespace, details.dependencyName, volume)
		source := filepath.Join(viper.GetString(config.ServicePathHost), path)
		volumes = append(volumes, mount.Mount{
			Source: source,
			Target: volume,
		})
		os.MkdirAll(filepath.Join(viper.GetString(config.ServicePathDocker), path), os.ModePerm)
	}
	for _, dep := range dependency.Volumesfrom {
		for _, volume := range service.Dependencies[dep].Volumes {
			path := filepath.Join(details.namespace, dep, volume)
			source := filepath.Join(viper.GetString(config.ServicePathHost), path)
			volumes = append(volumes, mount.Mount{
				Source: source,
				Target: volume,
			})
			os.MkdirAll(filepath.Join(viper.GetString(config.ServicePathDocker), path), os.ModePerm)
		}
	}
	return
}

func getDockerService(namespace string, dependencyName string) (dockerService swarm.Service, err error) {
	ctx := context.Background()
	cli, err := dockerCli()
	if err != nil {
		return
	}
	dockerServices, err := cli.ListServices(docker.ListServicesOptions{
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
	cli, err := dockerCli()
	if err != nil {
		return
	}
	containers, err := cli.ListContainers(docker.ListContainersOptions{
		Context: context.Background(),
		Limit:   1,
		Filters: map[string][]string{
			"name": []string{strings.Join([]string{namespace, dependencyName}, "_")},
			"label": []string{
				"com.docker.stack.namespace=" + namespace,
			},
		},
	})
	status = STOPPED
	if err == nil && len(containers) == 1 {
		status = RUNNING
	}
	return
}
