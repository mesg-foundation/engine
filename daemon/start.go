package daemon

import (
	"fmt"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	godocker "github.com/fsouza/go-dockerclient"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/docker"
	"github.com/spf13/viper"
)

func networkConfig() godocker.CreateNetworkOptions {
	return godocker.CreateNetworkOptions{
		Name:           sharedNetwork,
		CheckDuplicate: true, // Cannot have 2 network with the same name
		Driver:         "overlay",
		Labels: map[string]string{
			"com.docker.stack.namespace": sharedNetwork,
		},
	}
}

func serviceConfig(networkID string) godocker.CreateServiceOptions {
	return godocker.CreateServiceOptions{
		ServiceSpec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: name,
				Labels: map[string]string{
					"com.docker.stack.namespace": name,
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: &swarm.ContainerSpec{
					Image: image,
					Env: []string{
						"MESG.PATH=" + viper.GetString(config.MESGPath),
					},
					Mounts: []mount.Mount{
						mount.Mount{
							Source: socketPath,
							Target: socketPath,
						},
						mount.Mount{
							Source: viper.GetString(config.MESGPath),
							Target: "/mesg",
						},
					},
					Labels: map[string]string{
						"com.docker.stack.namespace": name,
					},
				},
			},
			Networks: []swarm.NetworkAttachmentConfig{
				swarm.NetworkAttachmentConfig{
					Target: sharedNetwork,
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
		},
	}
}

func Start() (err error) {
	client, err := docker.Client()
	if err != nil {
		return
	}

	network, err := docker.FindNetwork(sharedNetwork)
	if network == nil {
		fmt.Println("Create docker network")
		network, err = client.CreateNetwork(networkConfig())
		if err != nil {
			return
		}
	}

	fmt.Println("Create docker service")
	_, err = client.CreateService(serviceConfig("")) //network.ID))
	if err != nil {
		return
	}

	return
}
