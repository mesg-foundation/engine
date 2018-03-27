package cmdWorkflow

import (
	"fmt"
	"os"

	"github.com/mesg-foundation/application/cmd/utils"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/spf13/cobra"
)

// Pause run the pause command for a workflow
var Pause = &cobra.Command{
	Use:   "pause WORKFLOW_ID",
	Short: "Stop a workflow without killing it",
	Example: `mesg-cli workflow pause WORKFLOW_ID
mesg-cli workflow pause WORKFLOW_ID --account ACCOUNT_ID --confirm`,
	Run:               pauseHandler,
	DisableAutoGenTag: true,
}

func pauseHandler(cmd *cobra.Command, args []string) {
	account := cmdUtils.AccountFromFlagOrAsk(cmd, "Select an account:")
	var workflow = ""
	if len(args) > 0 {
		workflow = args[0]
	}
	if workflow == "" {
		// TODO add real list
		workflows := []string{"Workflow #1", "Workflow #2"}
		if survey.AskOne(&survey.Select{
			Message: "Choose the workflow to pause:",
			Default: workflows[0],
			Options: workflows,
		}, &workflow, nil) != nil {
			os.Exit(0)
		}
	}
	if !cmdUtils.Confirm(cmd, "Are you sure?") {
		return
	}
	// TODO pause the workflow onchain
	fmt.Println("Workflow pause called", args, account)
}

func init() {
	cmdUtils.Confirmable(Pause)
	cmdUtils.Accountable(Pause)
}
