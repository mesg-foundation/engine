package cmdWorkflow

import (
	"fmt"
	"time"

	"github.com/kyokomi/emoji"
	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
)

// Validate run the validate command for a workflow
var Validate = &cobra.Command{
	Use:               "validate FILE",
	Short:             "Validate a workflow",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli workflow validate workflow.yml",
	Run:               validateHandler,
	DisableAutoGenTag: true,
}

func validateHandler(cmd *cobra.Command, args []string) {
	s := cmdUtils.StartSpinner(cmdUtils.SpinnerOptions{Text: "Validation..."})
	time.Sleep(time.Second)
	s.Stop()
	fmt.Println(emoji.Sprint(":white_check_mark: Workflow valid"))
	// TODO validate the workflow
}
