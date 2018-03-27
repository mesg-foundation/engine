package cmdService

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"

	"github.com/spf13/cobra"
)

// List all the services
var List = &cobra.Command{
	Use:   "list",
	Short: "List all published services",
	Long: `This command returns all published services with basic information.
Optionally, you can filter the services published by a specific developer:
To have more details, see the [detail command](mesg-cli_service_detail.md).`,
	Example: `mesg-cli service list
mesg-cli service list --account ACCOUNT`,
	Run:               listHandler,
	DisableAutoGenTag: true,
}

func listHandler(cmd *cobra.Command, args []string) {
	// TODO Fetch details and display
	fmt.Println(cmdUtils.AccountFromFlagOrAsk(cmd, "Select an account:"))
	fmt.Println("- service1")
	fmt.Println("- service2")
	fmt.Println("- service3")
}

func init() {
	cmdUtils.Accountable(List)
}
