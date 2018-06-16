package cmd

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/mesg-foundation/core/container"
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
	// TODO: should improve this function with a waitFor
	status, err := daemon.Status()
	utils.HandleError(err)
	if status == container.RUNNING {
		fmt.Println(aurora.Green("MESG Core is running"))
	} else {
		fmt.Println(aurora.Brown("MESG Core is stopped"))
	}
}
