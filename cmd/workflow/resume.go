package cmdWorkflow

import (
	"fmt"
	"os"

	"github.com/mesg-foundation/application/cmd/utils"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/spf13/cobra"
)

// Resume run the resume command for a workflow
var Resume = &cobra.Command{
	Use:   "resume WORKFLOW_ID",
	Short: "Resume a previously paused workflow",
	Long: `Resume a previously paused workflow.

To pause a workflow, see the [pause command](mesg-cli_workflow_pause.md)`,
	Example: `mesg-cli workflow resume WORKFLOW_ID
mesg-cli workflow resume WORKFLOW_ID --account ACCOUNT_ID --confirm`,
	Run:               resumeHandler,
	DisableAutoGenTag: true,
}

func resumeHandler(cmd *cobra.Command, args []string) {
	account := cmdUtils.AccountFromFlagOrAsk(cmd, "Select an account:")
	var workflow = ""
	if len(args) > 0 {
		workflow = args[0]
	}
	if workflow == "" {
		// TODO add real list
		workflows := []string{"Workflow #1", "Workflow #2"}
		if survey.AskOne(&survey.Select{
			Message: "Choose the workflow to resume:",
			Default: workflows[0],
			Options: workflows,
		}, &workflow, nil) != nil {
			os.Exit(0)
		}
	}
	if !cmdUtils.Confirm(cmd, "Are you sure?") {
		return
	}
	// TODO resume the workflow
	fmt.Println("Workflow resume called", args, account)
}

func init() {
	cmdUtils.Confirmable(Resume)
	cmdUtils.Accountable(Resume)
}
