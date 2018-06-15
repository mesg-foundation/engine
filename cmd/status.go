package cmd

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/daemon"
	"github.com/spf13/cobra"
)

// Status command returns started services
var Status = &cobra.Command{
	Use:               "status",
	Short:             "Status of the MESG Core",
	Run:               statusHandler,
	DisableAutoGenTag: true,
}

func init() {
	RootCmd.AddCommand(Status)
}

func statusHandler(cmd *cobra.Command, args []string) {
	running, err := daemon.IsRunning()
	cmdUtils.HandleError(err)
	if running {
		fmt.Println(aurora.Green("MESG Core is running"))
	} else {
		fmt.Println(aurora.Brown("MESG Core is stopped"))
	}
}
