package container

import (
	"strings"
	"time"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
)

// ServiceOptions is a simplify version of swarm.ServiceSpec.
type ServiceOptions struct {
	Image           string
	Namespace       string
	Ports           []Port
	Mounts          []Mount
	Env             []string
	Args            []string
	Command         string
	Networks        []Network
	Labels          map[string]string
	StopGracePeriod *time.Duration
}

// Network keeps the network info for service.
type Network struct {
	// ID of the docker network.
	ID string

	// Alias is an optional attribute to name this service in the
	// network and be able to access to it using this name.
	Alias string
}

// Port is a simplify version of swarm.PortConfig.
type Port struct {
	Target    uint32
	Published uint32
}

// Mount is a simplify version of mount.Mount.
type Mount struct {
	Source string
	Target string
	Bind   bool
}

func (options *ServiceOptions) toSwarmServiceSpec(c *DockerContainer) swarm.ServiceSpec {
	namespace := c.namespace(options.Namespace)
	return swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: namespace,
			Labels: mergeStringMaps(options.Labels, map[string]string{
				"com.docker.stack.namespace": namespace,
				"com.docker.stack.image":     options.Image,
			}),
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image: options.Image,
				Labels: map[string]string{
					"com.docker.stack.namespace": namespace,
				},
				Env:             options.Env,
				Args:            options.Args,
				Command:         strings.Fields(options.Command),
				Mounts:          options.swarmMounts(),
				StopGracePeriod: options.StopGracePeriod,
			},
			Networks: options.swarmNetworks(),
		},
		EndpointSpec: &swarm.EndpointSpec{
			Ports: options.swarmPorts(),
		},
	}
}

func (options *ServiceOptions) swarmPorts() []swarm.PortConfig {
	ports := make([]swarm.PortConfig, len(options.Ports))
	for i, p := range options.Ports {
		ports[i] = swarm.PortConfig{
			Protocol:      swarm.PortConfigProtocolTCP,
			PublishMode:   swarm.PortConfigPublishModeIngress,
			TargetPort:    p.Target,
			PublishedPort: p.Published,
		}
	}
	return ports
}

func (options *ServiceOptions) swarmMounts() []mount.Mount {
	mounts := make([]mount.Mount, len(options.Mounts))
	for i, m := range options.Mounts {
		mountType := mount.TypeVolume
		if m.Bind {
			mountType = mount.TypeBind
		}
		mounts[i] = mount.Mount{
			Source: m.Source,
			Target: m.Target,
			Type:   mountType,
		}
	}
	return mounts
}

// swarmNetworks creates all necessary network attachment configurations for service.
// each network will be attached based on their networkID and an alias can be used to
// identify service in the network.
// aliases will make services accessible from other containers inside the same network.
func (options *ServiceOptions) swarmNetworks() []swarm.NetworkAttachmentConfig {
	networks := make([]swarm.NetworkAttachmentConfig, len(options.Networks))
	for i, network := range options.Networks {
		cfg := swarm.NetworkAttachmentConfig{
			Target: network.ID,
		}
		if network.Alias != "" {
			cfg.Aliases = []string{network.Alias}
		}
		networks[i] = cfg
	}
	return networks
}

func mergeStringMaps(m ...map[string]string) map[string]string {
	out := make(map[string]string)

	for i := range m {
		for k, v := range m[i] {
			out[k] = v
		}
	}
	return out
}
