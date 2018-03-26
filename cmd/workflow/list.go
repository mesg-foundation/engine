package cmdWorkflow

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
	fmt.Println("filter : ", cmd.Flag("account").Value.String())
	fmt.Println("- Workflow 1")
	fmt.Println("- Workflow 2")
	fmt.Println("- Workflow 3")
}

func init() {
	List.Flags().StringP("account", "a", "", "Filter workflows based on the account address")
}
