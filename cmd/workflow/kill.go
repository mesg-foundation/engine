package cmdWorkflow

import (
	"fmt"
	"time"

	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
)

// Kill a workflow
var Kill = &cobra.Command{
	Use:               "kill ID",
	Short:             "Kill a workflow and return the amount of money in this workflow",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli workflow kill xx",
	Run:               killHandler,
	DisableAutoGenTag: true,
}

func killHandler(cmd *cobra.Command, args []string) {
	account := cmdUtils.AccountFromFlagOrAsk(cmd, "Select an account")
	if !cmdUtils.Confirm(cmd, "Are you sure ?") {
		return
	}
	s := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Killing in progress..."})
	time.Sleep(2 * time.Second)
	s.Stop()
	// TODO kill the workflow
	fmt.Println("workflow killed", args, account)
}

func init() {
	cmdUtils.Confirmable(Kill)
	cmdUtils.Accountable(Kill)
}
