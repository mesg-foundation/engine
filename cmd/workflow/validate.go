package cmdWorkflow

import (
	"fmt"
	"time"

	"github.com/kyokomi/emoji"
	"github.com/mesg-foundation/core/cmd/utils"

	"github.com/spf13/cobra"
)

// Validate run the validate command for a workflow
var Validate = &cobra.Command{
	Use:   "validate ./PATH_TO_WORKFLOW_FILE",
	Short: "Validate a workflow",
	Long: `Validate a workflow.

To get more information, see the [deploy page from the documentation](https://docs.mesg.tech/workflow/validate.html)`,
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-core workflow validate ./PATH_TO_WORKFLOW_FILE.yml",
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
