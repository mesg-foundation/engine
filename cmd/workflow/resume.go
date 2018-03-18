package cmdWorkflow

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/spf13/cobra"
)

// Resume run the resume command for a workflow
var Resume = &cobra.Command{
	Use:               "resume ID",
	Short:             "Resume a paused workflow",
	Example:           "mesg-cli workflow resume xxx",
	Run:               resumeHandler,
	DisableAutoGenTag: true,
}

func resumeHandler(cmd *cobra.Command, args []string) {
	account := cmdUtils.AccountFromFlagOrAsk(cmd, "Select an account")
	var workflow = ""
	if len(args) > 0 {
		workflow = args[0]
	}
	if workflow == "" {
		// TODO add real list
		workflows := []string{"Workflow #1", "Workflow #2"}
		survey.AskOne(&survey.Select{
			Message: "Choose the workflow you want to resume",
			Default: workflows[0],
			Options: workflows,
		}, &workflow, nil)
	}
	if !cmdUtils.Confirm(cmd, "Are you sure ?") {
		return
	}
	// TODO resume the workflow
	fmt.Println("workflow resume called", args, account)
}

func init() {
	cmdUtils.Confirmable(Resume)
	cmdUtils.Accountable(Resume)
}
