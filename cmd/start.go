package cmd

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/daemon"
	"github.com/spf13/cobra"
)

// Start the MESG Core
var Start = &cobra.Command{
	Use:               "start",
	Short:             "Start the MESG Core",
	Run:               startHandler,
	DisableAutoGenTag: true,
}

func init() {
	RootCmd.AddCommand(Start)
}

func startHandler(cmd *cobra.Command, args []string) {
	running, err := daemon.IsRunning()
	cmdUtils.HandleError(err)
	if running {
		fmt.Println(aurora.Green("MESG Core is running"))
		return
	}
	cmdUtils.ShowSpinnerForFunc(cmdUtils.SpinnerOptions{Text: "Starting MESG Core..."}, func() {
		_, err = daemon.Start()
	})
	cmdUtils.HandleError(err)
	fmt.Println(aurora.Green("MESG Core is running"))
}
