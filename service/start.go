package service

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/swarm"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/daemon"
	"github.com/mesg-foundation/core/docker"
	"github.com/spf13/viper"
)

type dockerConfig struct {
	service    *Service
	dependency *Dependency
	name       string
	networkID  string
}

// Start a service
func (service *Service) Start() (dockerServices []*swarm.Service, err error) {
	status, err := service.Status()
	if err != nil || status == RUNNING {
		return
	}
	// If there is one but not all services running stop to restart all
	if status == PARTIAL {
		service.Stop()
	}
	network, err := docker.CreateNetwork([]string{service.Name}, "overlay")
	if err != nil {
		return
	}
	dockerServices = make([]*swarm.Service, len(service.Dependencies))
	i := 0
	for name, dependency := range service.Dependencies {
		dockerServices[i], err = startDocker(dockerConfig{
			service:    service,
			dependency: dependency,
			name:       name,
			networkID:  network.ID,
		})
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

// start will start a dependency container
func startDocker(c dockerConfig) (dockerService *swarm.Service, err error) {
	sharedNetwork, err := daemon.SharedNetwork()
	if err != nil {
		return
	}
	if sharedNetwork == nil {
		err = errors.New("Daemon shared network not found")
		return
	}
	return docker.StartService(&docker.ServiceOptions{
		Image:     c.dependency.Image,
		Namespace: []string{c.service.Name, c.name},
		Labels: map[string]string{
			dockerLabelServiceKey: c.service.Name,
		},
		Ports: c.dockerPorts(),
		Mounts: append(c.dockerVolumes(), docker.Mount{
			Source: viper.GetString(config.APIServiceSocketPath),
			Target: viper.GetString(config.APIServiceTargetPath),
		}),
		Env: []string{
			"MESG_ENDPOINT=" + viper.GetString(config.APIServiceTargetSocket),
			"MESG_ENDPOINT_TCP=" + docker.Namespace(daemon.Namespace()) + viper.GetString(config.APIClientTarget),
		},
		Args:       strings.Fields(c.dependency.Command),
		NetworksID: []string{c.networkID, sharedNetwork.ID},
	})
}

// dockerPorts extract ports from a Dependency and transform them to a swarm.PortConfig
func (c *dockerConfig) dockerPorts() (ports []docker.Port) {
	ports = make([]docker.Port, len(c.dependency.Ports))
	for i, p := range c.dependency.Ports {
		split := strings.Split(p, ":")
		published, _ := strconv.ParseUint(split[0], 10, 64)
		target := published
		if len(split) > 1 {
			target, _ = strconv.ParseUint(split[1], 10, 64)
		}
		ports[i] = docker.Port{
			Target:    uint32(target),
			Published: uint32(published),
		}
	}
	return
}

// dockerVolumes extract volumes from a Dependency and transform them to a Docker Mount
func (c *dockerConfig) dockerVolumes() (mounts []docker.Mount) {
	mounts = make([]docker.Mount, 0)
	for _, volume := range c.dependency.Volumes {
		path := filepath.Join(docker.Namespace([]string{c.service.Name, c.name}), volume)
		source := filepath.Join(viper.GetString(config.ServicePathHost), path)
		mounts = append(mounts, docker.Mount{
			Source: source,
			Target: volume,
		})
		os.MkdirAll(filepath.Join(viper.GetString(config.ServicePathDocker), path), os.ModePerm)
	}
	for _, depString := range c.dependency.Volumesfrom {
		for _, volume := range c.service.Dependencies[depString].Volumes {
			path := filepath.Join(docker.Namespace([]string{c.service.Name, depString}), volume)
			source := filepath.Join(viper.GetString(config.ServicePathHost), path)
			mounts = append(mounts, docker.Mount{
				Source: source,
				Target: volume,
			})
			os.MkdirAll(filepath.Join(viper.GetString(config.ServicePathDocker), path), os.ModePerm)
		}
	}
	return
}
