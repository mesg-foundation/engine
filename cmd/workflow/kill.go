package cmdWorkflow

import (
	"fmt"
	"time"

	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
)

// Kill a workflow
var Kill = &cobra.Command{
	Use:   "kill WORKFLOW_ID",
	Short: "Kill a workflow from the Network and get back its token",
	Args:  cobra.MinimumNArgs(1),
	Example: `mesg-cli workflow kill WORKFLOW_ID
mesg-cli workflow kill WORKFLOW_ID --account ACCOUNT --confirm`,
	Run:               killHandler,
	DisableAutoGenTag: true,
}

func killHandler(cmd *cobra.Command, args []string) {
	account := cmdUtils.AccountFromFlagOrAsk(cmd, "Select an account:")
	if !cmdUtils.Confirm(cmd, "Are you sure?") {
		return
	}
	s := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Killing in progress..."})
	time.Sleep(2 * time.Second)
	s.Stop()
	// TODO kill the workflow
	fmt.Println("Workflow killed with success", args, account)
}

func init() {
	cmdUtils.Confirmable(Kill)
	cmdUtils.Accountable(Kill)
}
