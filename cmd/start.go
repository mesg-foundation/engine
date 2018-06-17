package cmd

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/container"
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
	status, err := daemon.Status()
	utils.HandleError(err)
	if status == container.RUNNING {
		fmt.Println(aurora.Green("MESG Core is running"))
		return
	}
	utils.ShowSpinnerForFunc(utils.SpinnerOptions{Text: "Starting MESG Core..."}, func() {
		_, err = daemon.Start()
	})
	utils.HandleError(err)
	fmt.Println(aurora.Green("MESG Core is running"))
}
