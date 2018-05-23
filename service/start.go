package service

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	godocker "github.com/fsouza/go-dockerclient"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/daemon"
	"github.com/mesg-foundation/core/docker"
	"github.com/spf13/viper"
)

type dockerConfig struct {
	service    *Service
	dependency *Dependency
	name       string
}

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
		dockerServices[i], err = startDocker(dockerConfig{
			service:    service,
			dependency: dependency,
			name:       name,
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
	daemonIP, err := daemon.IP()
	if err != nil {
		return
	}
	sharedNetworkID, err := daemon.SharedNetworkID()
	if err != nil {
		return
	}
	return docker.StartService(godocker.CreateServiceOptions{
		ServiceSpec: swarm.ServiceSpec{
			Annotations: c.dockerAnnotations(),
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: c.dockerContainerSpec(daemonIP),
			},
			EndpointSpec: &swarm.EndpointSpec{
				Ports: c.dockerPorts(),
			},
			Networks: []swarm.NetworkAttachmentConfig{
				swarm.NetworkAttachmentConfig{
					Target: sharedNetworkID,
				},
			},
		},
	})
}

// dockerPorts extract ports from a Dependency and transform them to a swarm.PortConfig
func (c *dockerConfig) dockerPorts() (ports []swarm.PortConfig) {
	ports = make([]swarm.PortConfig, len(c.dependency.Ports))
	for i, p := range c.dependency.Ports {
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

// dockerVolumes extract volumes from a Dependency and transform them to a Docker Mount
func (c *dockerConfig) dockerVolumes() (volumes []mount.Mount) {
	volumes = make([]mount.Mount, 0)
	for _, volume := range c.dependency.Volumes {
		path := filepath.Join(docker.Namespace([]string{c.service.Name, c.name}), volume)
		source := filepath.Join(viper.GetString(config.ServicePathHost), path)
		volumes = append(volumes, mount.Mount{
			Source: source,
			Target: volume,
		})
		os.MkdirAll(filepath.Join(viper.GetString(config.ServicePathDocker), path), os.ModePerm)
	}
	for _, depString := range c.dependency.Volumesfrom {
		for _, volume := range c.service.Dependencies[depString].Volumes {
			path := filepath.Join(docker.Namespace([]string{c.service.Name, depString}), volume)
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

func (c *dockerConfig) dockerAnnotations() swarm.Annotations {
	namespace := docker.Namespace([]string{c.service.Name, c.name})
	return swarm.Annotations{
		Name: namespace,
		Labels: map[string]string{
			"com.docker.stack.image":     c.dependency.Image,
			"com.docker.stack.namespace": namespace,
			"mesg.service":               c.service.Name,
		},
	}
}

func (c *dockerConfig) dockerContainerSpec(daemonIP string) *swarm.ContainerSpec {
	namespace := docker.Namespace([]string{c.service.Name, c.name})
	return &swarm.ContainerSpec{
		Image: c.dependency.Image,
		Args:  strings.Fields(c.dependency.Command),
		Env: []string{
			"MESG_ENDPOINT=" + viper.GetString(config.APIServiceTargetSocket),
			"MESG_ENDPOINT_TCP=" + daemonIP + "" + viper.GetString(config.APIClientTarget),
		},
		Mounts: append(c.dockerVolumes(), mount.Mount{
			Source: viper.GetString(config.APIServiceSocketPath),
			Target: viper.GetString(config.APIServiceTargetPath),
		}),
		Labels: map[string]string{
			"com.docker.stack.namespace": namespace,
		},
	}
}
