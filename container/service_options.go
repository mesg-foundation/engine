package container

import (
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
)

// ServiceOptions is a simplify version of swarm.ServiceSpec that can be merge to it
type ServiceOptions struct {
	Image       string
	Namespace   []string
	Ports       []Port
	Mounts      []Mount
	Env         []string
	Args        []string
	NetworksID  []string
	Labels      map[string]string
	ServiceSpec swarm.ServiceSpec
}

// Port is a simplify version of swarm.PortConfig
type Port struct {
	Target    uint32
	Published uint32
}

// Mount is a simplify version of mount.Mount
type Mount struct {
	Source string
	Target string
}

func (options *ServiceOptions) merge() {
	options.initCreateServiceOptions()
	options.mergeNamespace()
	options.mergeImage()
	options.mergeLabels()
	options.mergePorts()
	options.mergeMounts()
	options.mergeEnv()
	options.mergeArgs()
	options.mergeNetworks()
}

func (options *ServiceOptions) initCreateServiceOptions() {
	if options.ServiceSpec.Annotations.Labels == nil {
		options.ServiceSpec.Annotations.Labels = make(map[string]string)
	}
	if options.ServiceSpec.EndpointSpec == nil {
		options.ServiceSpec.EndpointSpec = &swarm.EndpointSpec{}
	}
	if options.ServiceSpec.EndpointSpec.Ports == nil {
		options.ServiceSpec.EndpointSpec.Ports = make([]swarm.PortConfig, 0)
	}
	if options.ServiceSpec.TaskTemplate.ContainerSpec == nil {
		options.ServiceSpec.TaskTemplate.ContainerSpec = &swarm.ContainerSpec{}
	}
	if options.ServiceSpec.TaskTemplate.Networks == nil {
		options.ServiceSpec.TaskTemplate.Networks = make([]swarm.NetworkAttachmentConfig, 0)
	}
	if options.ServiceSpec.TaskTemplate.ContainerSpec.Args == nil {
		options.ServiceSpec.TaskTemplate.ContainerSpec.Args = make([]string, 0)
	}
	if options.ServiceSpec.TaskTemplate.ContainerSpec.Env == nil {
		options.ServiceSpec.TaskTemplate.ContainerSpec.Env = make([]string, 0)
	}
	if options.ServiceSpec.TaskTemplate.ContainerSpec.Mounts == nil {
		options.ServiceSpec.TaskTemplate.ContainerSpec.Mounts = make([]mount.Mount, 0)
	}
	if options.ServiceSpec.TaskTemplate.ContainerSpec.Labels == nil {
		options.ServiceSpec.TaskTemplate.ContainerSpec.Labels = make(map[string]string)
	}
}

func (options *ServiceOptions) mergeNamespace() {
	namespace := Namespace(options.Namespace)
	options.ServiceSpec.Annotations.Name = namespace
	options.ServiceSpec.Annotations.Labels["com.docker.stack.namespace"] = namespace
	options.ServiceSpec.TaskTemplate.ContainerSpec.Labels["com.docker.stack.namespace"] = namespace
}

func (options *ServiceOptions) mergeImage() {
	options.ServiceSpec.Annotations.Labels["com.docker.stack.image"] = options.Image
	options.ServiceSpec.TaskTemplate.ContainerSpec.Image = options.Image
}

func (options *ServiceOptions) mergeLabels() {
	for k, v := range options.Labels {
		options.ServiceSpec.Annotations.Labels[k] = v
	}
}

func (options *ServiceOptions) mergePorts() {
	ports := make([]swarm.PortConfig, len(options.Ports))
	for i, p := range options.Ports {
		ports[i] = swarm.PortConfig{
			Protocol:      swarm.PortConfigProtocolTCP,
			PublishMode:   swarm.PortConfigPublishModeIngress,
			TargetPort:    p.Target,
			PublishedPort: p.Published,
		}
	}
	options.ServiceSpec.EndpointSpec.Ports = append(options.ServiceSpec.EndpointSpec.Ports, ports...)
}

func (options *ServiceOptions) mergeMounts() {
	mounts := make([]mount.Mount, len(options.Mounts))
	for i, m := range options.Mounts {
		mounts[i] = mount.Mount{
			Source: m.Source,
			Target: m.Target,
		}
	}
	options.ServiceSpec.TaskTemplate.ContainerSpec.Mounts = append(options.ServiceSpec.TaskTemplate.ContainerSpec.Mounts, mounts...)
}

func (options *ServiceOptions) mergeEnv() {
	options.ServiceSpec.TaskTemplate.ContainerSpec.Env = append(options.ServiceSpec.TaskTemplate.ContainerSpec.Env, options.Env...)
}

func (options *ServiceOptions) mergeArgs() {
	options.ServiceSpec.TaskTemplate.ContainerSpec.Args = append(options.ServiceSpec.TaskTemplate.ContainerSpec.Args, options.Args...)
}

func (options *ServiceOptions) mergeNetworks() {
	for _, networkID := range options.NetworksID {
		options.ServiceSpec.TaskTemplate.Networks = append(options.ServiceSpec.TaskTemplate.Networks, swarm.NetworkAttachmentConfig{
			Target: networkID,
		})
	}
}
