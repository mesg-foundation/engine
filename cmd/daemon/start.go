package daemon

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/daemon"
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

		spinner := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Starting the daemon"})
		err = <-daemon.WaitForContainerToRun()
		spinner.Stop()
		if err != nil {
			fmt.Println(aurora.Red(err))
			return
		}
	}

	fmt.Println(aurora.Green("Daemon is running"))
}
