package cmdWorkflow

import (
	"fmt"

	"github.com/spf13/cobra"
)

// List workflows
var List = &cobra.Command{
	Use:   "list",
	Short: "List all deployed workflows",
	Long: `List all workflows deployed on the Network.

Optionally, you can filter the workflows deployed by a specific account.

This command will return basic information. To have more details, see the [detail command](mesg-core_workflow_detail.md).`,
	Example: `mesg-core workflow list
mesg-core workflow list --account ACCOUNT`,
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
	List.Flags().StringP("account", "a", "", "Filter workflows by a specific account")
}
