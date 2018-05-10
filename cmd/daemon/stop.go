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
	service, err := service()
	if err != nil {
		fmt.Println(aurora.Red(err))
	}
	if service != nil {
		client, err := docker.NewClientFromEnv()
		if err != nil {
			fmt.Println(aurora.Red(err))
		}
		err = client.RemoveService(docker.RemoveServiceOptions{
			Context: context.Background(),
			ID:      service.ID,
		})
		if err != nil {
			fmt.Println(aurora.Red(err))
		}
	}
	fmt.Println(aurora.Green("Daemon stopped"))
}
