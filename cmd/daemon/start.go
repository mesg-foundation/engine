package daemon

import (
	"context"
	"fmt"
	"log"
	"time"

	godocker "github.com/fsouza/go-dockerclient"
	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/daemon"
	"github.com/mesg-foundation/core/docker"
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
	running, err := daemon.IsRunning()
	if err != nil {
		fmt.Println(aurora.Red(err))
		return
	}
	if !running {
		_, err = daemon.Start()
		if err != nil {
			fmt.Println(aurora.Red(err))
			return
		}

		// spinner := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Starting the daemon"})
		go func() {
			client, _ := docker.Client()
			for {
				containers, _ := client.ListContainers(godocker.ListContainersOptions{
					Context: context.Background(),
				})
				for _, c := range containers {
					fmt.Println("id:", c.ID, "image:", c.Image, "name:", c.Names, "state", c.State, "status", c.Status)
					log.Println("id:", c.ID, "image:", c.Image, "name:", c.Names, "state", c.State, "status", c.Status)
				}
				time.Sleep(1 * time.Second)
			}
		}()
		err = <-daemon.WaitForContainerToRun()
		// spinner.Stop()
		if err != nil {
			fmt.Println(aurora.Red(err))
			return
		}
	}

	fmt.Println(aurora.Green("Daemon is running"))
}
