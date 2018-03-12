package cmdWorkflow

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"
	survey "gopkg.in/AlecAivazis/survey.v1"

	"github.com/spf13/cobra"
)

// Pause run the pause command for a workflow
var Pause = &cobra.Command{
	Use:               "pause",
	Short:             "Pause a workflow",
	Example:           "mesg-cli workflow pause xxx",
	Run:               pauseHandler,
	DisableAutoGenTag: true,
}

func pauseHandler(cmd *cobra.Command, args []string) {
	var workflow = ""
	if len(args) > 0 {
		workflow = args[0]
	}
	if workflow == "" {
		// TODO add real list
		workflows := []string{"Workflow #1", "Workflow #2"}
		survey.AskOne(&survey.Select{
			Message: "Choose the workflow you want to pause",
			Default: workflows[0],
			Options: workflows,
		}, &workflow, nil)
	}
	if !cmdUtils.Confirm(cmd, "Are you sure ?") {
		return
	}
	// TODO pause the workflow onchain
	fmt.Println("workflow pause called", args)
}

func init() {
	Pause.Flags().BoolP("confirm", "c", false, "Confirm")
}
