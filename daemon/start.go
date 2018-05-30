package daemon

import (
	"path/filepath"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/spf13/viper"
)

// Start the docker daemon
func Start() (serviceID string, err error) {
	running, err := IsRunning()
	if err != nil {
		return
	}
	if running == true {
		return
	}
	sharedNetworkID, err := container.SharedNetworkID()
	if err != nil {
		return
	}
	return container.StartService(serviceSpec(sharedNetworkID))
}

func serviceSpec(networkID string) swarm.ServiceSpec {
	namespace := container.Namespace([]string{name})
	return swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: namespace,
			Labels: map[string]string{
				"com.docker.stack.image":     viper.GetString(config.DaemonImage),
				"com.docker.stack.namespace": namespace,
			},
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image: viper.GetString(config.DaemonImage),
				Env: []string{
					"MESG.PATH=/mesg",
					"API.SERVICE.SOCKETPATH=" + filepath.Join(viper.GetString(config.MESGPath), "server.sock"),
				},
				Mounts: []mount.Mount{
					mount.Mount{
						Source: dockerSocket,
						Target: dockerSocket,
					},
					mount.Mount{
						Source: viper.GetString(config.MESGPath),
						Target: "/mesg",
					},
				},
				Labels: map[string]string{
					"com.docker.stack.namespace": namespace,
				},
			},
		},
		EndpointSpec: &swarm.EndpointSpec{
			Ports: []swarm.PortConfig{
				swarm.PortConfig{
					Protocol:      swarm.PortConfigProtocolTCP,
					PublishMode:   swarm.PortConfigPublishModeIngress,
					TargetPort:    50052,
					PublishedPort: 50052,
				},
			},
		},
		Networks: []swarm.NetworkAttachmentConfig{
			swarm.NetworkAttachmentConfig{
				Target: networkID,
			},
		},
	}
}
