package daemon

import (
	"github.com/docker/docker/api/types/swarm"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/docker"
	"github.com/spf13/viper"
)

// Start the docker daemon
func Start() (service *swarm.Service, err error) {
	network, err := docker.CreateNetwork(sharedNetwork)
	if err != nil {
		return
	}
	return docker.StartService(&docker.ServiceOptions{
		Image:     image,
		Namespace: []string{name},
		Ports: []docker.Port{
			docker.Port{
				Target:    50052,
				Published: 50052,
			},
		},
		Mounts: []docker.Mount{
			docker.Mount{
				Source: dockerSocket,
				Target: dockerSocket,
			},
			docker.Mount{
				Source: viper.GetString(config.MESGPath),
				Target: "/mesg",
			},
		},
		// Env: []string{
		// "MESG.PATH=" + viper.GetString(config.MESGPath), //TODO: is it really useful?
		// },
		NetworksID: []string{network.ID},
		// Service: &godocker.CreateServiceOptions{
		// ServiceSpec: swarm.ServiceSpec{
		// Annotations: swarm.Annotations{
		// Name: namespace,
		// Labels: map[string]string{
		// "com.docker.stack.image": image,
		// "com.docker.stack.namespace": namespace,
		// },
		// },
		// TaskTemplate: swarm.TaskSpec{
		// 	ContainerSpec: &swarm.ContainerSpec{
		// Image: image,
		// Env: []string{
		// 	"MESG.PATH=" + viper.GetString(config.MESGPath),
		// },
		// Mounts: []mount.Mount{
		// 	mount.Mount{
		// 		Source: dockerSocket,
		// 		Target: dockerSocket,
		// 	},
		// 	mount.Mount{
		// 		Source: viper.GetString(config.MESGPath),
		// 		Target: "/mesg",
		// 	},
		// },
		// Labels: map[string]string{
		// "com.docker.stack.namespace": namespace,
		// },
		// 	},
		// },
		// Networks: []swarm.NetworkAttachmentConfig{
		// 	swarm.NetworkAttachmentConfig{
		// 		Target: sharedNetwork,
		// 	},
		// },
		// EndpointSpec: &swarm.EndpointSpec{
		// 	Ports: []swarm.PortConfig{
		// 		swarm.PortConfig{
		// 			Protocol:      swarm.PortConfigProtocolTCP,
		// 			PublishMode:   swarm.PortConfigPublishModeIngress,
		// 			TargetPort:    50052,
		// 			PublishedPort: 50052,
		// 		},
		// 	},
		// },
		// },
		// },
	})
}
