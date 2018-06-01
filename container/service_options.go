package container

import (
	"os"
	"strconv"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
)

// ServiceOptions is a simplify version of swarm.ServiceSpec that can be created it.
type ServiceOptions struct {
	Image      string
	Namespace  []string
	Ports      []Port
	Mounts     []Mount
	Env        []string
	Args       []string
	NetworksID []string
	Labels     map[string]string
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

func (options *ServiceOptions) toSwarmServiceSpec() (service swarm.ServiceSpec) {
	namespace := Namespace(options.Namespace)
	service = swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: namespace,
			Labels: mergeLabels(options.Labels, map[string]string{
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
				Env:    options.Env,
				Args:   options.Args,
				Mounts: options.swarmMounts(),
			},
			Networks: options.swarmNetworks(),
		},
		EndpointSpec: &swarm.EndpointSpec{
			Ports: options.swarmPorts(),
		},
	}
	return
}

func (options *ServiceOptions) swarmPorts() (ports []swarm.PortConfig) {
	ports = make([]swarm.PortConfig, len(options.Ports))
	for i, p := range options.Ports {
		ports[i] = swarm.PortConfig{
			Protocol:      swarm.PortConfigProtocolTCP,
			PublishMode:   swarm.PortConfigPublishModeIngress,
			TargetPort:    p.Target,
			PublishedPort: p.Published,
		}
	}
	return
}

func (options *ServiceOptions) swarmMounts() (mounts []mount.Mount) {
	// hack for preventing mount when in CircleCI
	circleCI, errCircle := strconv.ParseBool(os.Getenv("CIRCLECI"))
	if errCircle == nil && circleCI {
		return
	}
	mounts = make([]mount.Mount, len(options.Mounts))
	for i, m := range options.Mounts {
		mounts[i] = mount.Mount{
			Source: m.Source,
			Target: m.Target,
		}
	}
	return
}

func (options *ServiceOptions) swarmNetworks() (networks []swarm.NetworkAttachmentConfig) {
	networks = make([]swarm.NetworkAttachmentConfig, len(options.NetworksID))
	for i, networkID := range options.NetworksID {
		networks[i] = swarm.NetworkAttachmentConfig{
			Target: networkID,
		}
	}
	return
}

func mergeLabels(l1 map[string]string, l2 map[string]string) map[string]string {
	if l1 == nil {
		l1 = make(map[string]string)
	}
	for k, v := range l2 {
		l1[k] = v
	}
	return l1
}
