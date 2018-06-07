package cmd

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
	Short:             "Start the core",
	Run:               startHandler,
	DisableAutoGenTag: true,
}

func init() {
	RootCmd.AddCommand(Start)
}

func startHandler(cmd *cobra.Command, args []string) {
	running, err := daemon.IsRunning()
	if err != nil {
		fmt.Println(aurora.Red(err))
		return
	}
	if running {
		fmt.Println(aurora.Green("Core is running"))
		return
	}
	cmdUtils.ShowSpinnerForFunc(cmdUtils.SpinnerOptions{Text: "Starting core..."}, func() {
		_, err = daemon.Start()
	})
	if err != nil {
		fmt.Println(aurora.Red(err))
		return
	}
	fmt.Println(aurora.Green("Core is running"))
}
