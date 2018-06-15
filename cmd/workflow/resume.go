package workflow

import (
	"fmt"
	"os"

	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/spf13/cobra"
)

// Resume run the resume command for a workflow
var Resume = &cobra.Command{
	Use:   "resume WORKFLOW_ID",
	Short: "Resume a previously paused workflow",
	Long: `Resume a previously paused workflow.

To pause a workflow, see the [pause command](mesg-core_workflow_pause.md)`,
	Example: `mesg-core workflow resume WORKFLOW_ID
mesg-core workflow resume WORKFLOW_ID --account ACCOUNT_ID --confirm`,
	Run:               resumeHandler,
	DisableAutoGenTag: true,
}

func resumeHandler(cmd *cobra.Command, args []string) {
	account := utils.AccountFromFlagOrAsk(cmd, "Select an account:")
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
	if !utils.Confirm(cmd, "Are you sure?") {
		return
	}
	// TODO resume the workflow
	fmt.Println("Workflow resume called", args, account)
}

func init() {
	utils.Confirmable(Resume)
	utils.Accountable(Resume)
}
