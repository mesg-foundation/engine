package daemon

import (
	"fmt"
	"time"

	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/service"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/mesg-foundation/core/config"
	"github.com/spf13/viper"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// Start the daemon
var Start = &cobra.Command{
	Use:               "start",
	Short:             "Start the daemon",
	Run:               startHandler,
	DisableAutoGenTag: true,
}

func startHandler(cmd *cobra.Command, args []string) {
	running, err := isRunning()
	if err != nil {
		fmt.Println(aurora.Red(err))
		return
	}
	if !running {
		client, err := service.DockerClient()
		if err != nil {
			fmt.Println(aurora.Red(err))
			return
		}

		network, err := service.SharedNetwork(client)
		if err != nil {
			fmt.Println(aurora.Red(err))
			return
		}

		_, err = client.CreateService(serviceConfig(&network))
		if err != nil {
			fmt.Println(aurora.Red(err))
			return
		}

		spinner := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Starting the daemon"})
		for {
			time.Sleep(500 * time.Millisecond)
			container, _ := getContainer()
			if container != nil {
				break
			}
		}
		spinner.Stop()
	}

	fmt.Println(aurora.Green("Daemon is running"))
}

func serviceConfig(network *docker.Network) docker.CreateServiceOptions {
	return docker.CreateServiceOptions{
		ServiceSpec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: name,
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: &swarm.ContainerSpec{
					Image: viper.GetString(config.DaemonImage),
					Env: []string{
						"MESG.PATH=/mesg",
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
					Target: network.ID,
				},
			},
		},
	}
}
