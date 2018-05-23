package docker

import (
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	godocker "github.com/fsouza/go-dockerclient"
)

func (options *ServiceOptions) merge() {
	if options.CreateServiceOptions == nil {
		options.CreateServiceOptions = &godocker.CreateServiceOptions{}
	}
	service := options.CreateServiceOptions
	if service.TaskTemplate.ContainerSpec == nil {
		service.TaskTemplate.ContainerSpec = &swarm.ContainerSpec{}
	}
	if service.EndpointSpec == nil {
		service.EndpointSpec = &swarm.EndpointSpec{}
	}
	if service.Annotations.Labels == nil {
		service.Annotations.Labels = make(map[string]string)
	}
	if service.TaskTemplate.ContainerSpec.Labels == nil {
		service.TaskTemplate.ContainerSpec.Labels = make(map[string]string)
	}
	if service.EndpointSpec.Ports == nil {
		service.EndpointSpec.Ports = make([]swarm.PortConfig, 0)
	}
	if service.TaskTemplate.ContainerSpec.Mounts == nil {
		service.TaskTemplate.ContainerSpec.Mounts = make([]mount.Mount, 0)
	}
	if service.TaskTemplate.ContainerSpec.Env == nil {
		service.TaskTemplate.ContainerSpec.Env = make([]string, 0)
	}
	if service.TaskTemplate.ContainerSpec.Args == nil {
		service.TaskTemplate.ContainerSpec.Args = make([]string, 0)
	}
	if service.Networks == nil {
		service.Networks = make([]swarm.NetworkAttachmentConfig, 0)
	}

	options.mergeNamespace()
	options.mergeImage()
	options.mergeLabels()
	options.mergePorts()
	options.mergeMounts()
	options.mergeEnv()
	options.mergeArgs()
	options.mergeNetworks()
}

func (options *ServiceOptions) mergeNamespace() {
	service := options.CreateServiceOptions
	namespace := Namespace(options.Namespace)
	service.Annotations.Name = namespace
	service.Annotations.Labels["com.docker.stack.namespace"] = namespace
	service.TaskTemplate.ContainerSpec.Labels["com.docker.stack.namespace"] = namespace
}

func (options *ServiceOptions) mergeImage() {
	service := options.CreateServiceOptions
	service.Annotations.Labels["com.docker.stack.image"] = options.Image
	service.TaskTemplate.ContainerSpec.Image = options.Image
}

func (options *ServiceOptions) mergeLabels() {
	service := options.CreateServiceOptions
	for k, v := range options.Labels {
		service.Annotations.Labels[k] = v
	}
}

func (options *ServiceOptions) mergePorts() {
	service := options.CreateServiceOptions
	ports := make([]swarm.PortConfig, len(options.Ports))
	for i, p := range options.Ports {
		ports[i] = swarm.PortConfig{
			Protocol:      swarm.PortConfigProtocolTCP,
			PublishMode:   swarm.PortConfigPublishModeIngress,
			TargetPort:    p.Target,
			PublishedPort: p.Published,
		}
	}
	service.EndpointSpec.Ports = append(service.EndpointSpec.Ports, ports...)
}

func (options *ServiceOptions) mergeMounts() {
	service := options.CreateServiceOptions
	mounts := make([]mount.Mount, len(options.Mounts))
	for i, m := range options.Mounts {
		mounts[i] = mount.Mount{
			Source: m.Source,
			Target: m.Target,
		}
	}
	service.TaskTemplate.ContainerSpec.Mounts = append(service.TaskTemplate.ContainerSpec.Mounts, mounts...)
}

func (options *ServiceOptions) mergeEnv() {
	service := options.CreateServiceOptions
	service.TaskTemplate.ContainerSpec.Env = append(service.TaskTemplate.ContainerSpec.Env, options.Env...)
}

func (options *ServiceOptions) mergeArgs() {
	service := options.CreateServiceOptions
	service.TaskTemplate.ContainerSpec.Args = append(service.TaskTemplate.ContainerSpec.Args, options.Args...)
}

func (options *ServiceOptions) mergeNetworks() {
	service := options.CreateServiceOptions
	for _, networkID := range options.NetworksID {
		service.Networks = append(service.Networks, swarm.NetworkAttachmentConfig{
			Target: networkID,
		})
	}
}
