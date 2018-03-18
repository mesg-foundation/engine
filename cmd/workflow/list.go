package cmdWorkflow

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
)

// List all the workflows
var List = &cobra.Command{
	Use:               "list",
	Short:             "List of workflows that an account already deployed on the Network",
	Example:           "mesg-cli workflow list",
	Run:               listHandler,
	DisableAutoGenTag: true,
}

func listHandler(cmd *cobra.Command, args []string) {
	// TODO Fetch details and display
	fmt.Println(cmdUtils.AccountFromFlagOrAsk(cmd, "Select an account"))
	fmt.Println("- workflow1")
	fmt.Println("- workflow2")
	fmt.Println("- workflow3")
}

func init() {
	cmdUtils.Accountable(List)
}
