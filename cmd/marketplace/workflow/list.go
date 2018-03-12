package cmdWorkflowMarketPlace

import (
	"fmt"

	"github.com/spf13/cobra"
)

// List workflows
var List = &cobra.Command{
	Use:               "list",
	Short:             "List all workflows of the marketplace",
	Example:           "mesg-cli marketplace service list",
	Run:               listHandler,
	DisableAutoGenTag: true,
}

func listHandler(cmd *cobra.Command, args []string) {
	// TODO list
	fmt.Println("...")
}
