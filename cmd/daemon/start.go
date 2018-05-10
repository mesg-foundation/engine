package daemon

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/mesg-foundation/core/cmd/utils"

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
		client, err := docker.NewClientFromEnv()
		if err != nil {
			fmt.Println(aurora.Red(err))
			return
		}

		fmt.Println("Download docker image")
		var stream bytes.Buffer
		go func() {
			err = client.PullImage(docker.PullImageOptions{
				OutputStream:  &stream,
				RawJSONStream: false,
				Repository:    image,
			}, docker.AuthConfiguration{})
			if err != nil {
				fmt.Println(aurora.Red(err))
				os.Exit(1)
			}
		}()
		buf := make([]byte, 1024)
		for {
			n, err := stream.Read(buf)
			if err != nil {
				fmt.Println(aurora.Red(err))
				// os.Exit(1)
				// break
			}
			if n != 0 {
				fmt.Print(string(buf[:n]))
			}
			time.Sleep(500 * time.Millisecond)
		}

		_, err = client.CreateService(serviceConfig())
		if err != nil {
			fmt.Println(aurora.Red(err))
			return
		}

		spinner := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Starting the daemon"})
		for {
			time.Sleep(500 * time.Millisecond)
			container, _ := container()
			if container != nil {
				break
			}
		}
		spinner.Stop()
	}

	fmt.Println(aurora.Green("Daemon is running"))
}

func serviceConfig() docker.CreateServiceOptions {
	return docker.CreateServiceOptions{
		ServiceSpec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: name,
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
