package cmdWorkflow

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Detail of the workflow
var Detail = &cobra.Command{
	Use:               "detail ID",
	Short:             "List all details of a workflow",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli marketplace workflow detail XX",
	Run:               detailHandler,
	DisableAutoGenTag: true,
}

func detailHandler(cmd *cobra.Command, args []string) {
	// TODO Details workflow
	fmt.Println(args)
}
