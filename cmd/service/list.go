package cmdService

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
)

// List all the services
var List = &cobra.Command{
	Use:               "list",
	Short:             "List of services that a account already deployed on the Network",
	Example:           "mesg-cli service list",
	Run:               listHandler,
	DisableAutoGenTag: true,
}

func listHandler(cmd *cobra.Command, args []string) {
	// TODO Fetch details and display
	fmt.Println(cmdUtils.AccountFromFlagOrAsk(cmd, "Select an account"))
	fmt.Println("- service1")
	fmt.Println("- service2")
	fmt.Println("- service3")
}

func init() {
	cmdUtils.Accountable(List)
}
