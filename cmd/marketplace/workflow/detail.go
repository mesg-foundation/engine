package cmdWorkflowMarketPlace

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stake float64
var duration int

// Detail of the workflow
var Detail = &cobra.Command{
	Use:               "detail ID",
	Short:             "Details of a workflow",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli marketplace workflow detail XX",
	Run:               detailHandler,
	DisableAutoGenTag: true,
}

func detailHandler(cmd *cobra.Command, args []string) {
	// TODO Details workflow
	fmt.Println(args)
}
