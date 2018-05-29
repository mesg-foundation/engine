package daemon

import (
	"context"
	"fmt"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// Stop the daemon
var Stop = &cobra.Command{
	Use:               "stop",
	Short:             "Stop the daemon",
	Run:               stopHandler,
	DisableAutoGenTag: true,
}

func stopHandler(cmd *cobra.Command, args []string) {
	service, err := getService()
	if err != nil {
		fmt.Println(aurora.Red(err))
		return
	}
	if service != nil {
		client, err := docker.NewClientFromEnv()
		if err != nil {
			fmt.Println(aurora.Red(err))
			return
		}
		err = client.RemoveService(docker.RemoveServiceOptions{
			Context: context.Background(),
			ID:      service.ID,
		})
		if err != nil {
			fmt.Println(aurora.Red(err))
			return
		}
	}
	fmt.Println(aurora.Green("Daemon stopped"))
}
