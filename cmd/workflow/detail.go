package cmdWorkflow

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Detail of the workflow
var Detail = &cobra.Command{
	Use:               "detail WORKFLOW_ID",
	Short:             "List details of a workflow",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-core workflow detail WORKFLOW_ID",
	Run:               detailHandler,
	DisableAutoGenTag: true,
}

func detailHandler(cmd *cobra.Command, args []string) {
	// TODO Details workflow
	fmt.Println(args)
}
