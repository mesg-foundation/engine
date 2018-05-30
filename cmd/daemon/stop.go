package daemon

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/container"
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
		err = container.StopService([]string{name})
		if err != nil {
			fmt.Println(aurora.Red(err))
			return
		}
	}
	fmt.Println(aurora.Green("Daemon stopped"))
}
