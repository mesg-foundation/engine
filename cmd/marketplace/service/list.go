package cmdServiceMarketPlace

import (
	"fmt"

	"github.com/spf13/cobra"
)

// List all the services
var List = &cobra.Command{
	Use:               "list SERVICE",
	Short:             "Show the list of services",
	Long:              "List all published services. Optionally, filter for a specific account.",
	Example:           "mesg-cli marketplace service list",
	Run:               listHandler,
	DisableAutoGenTag: true,
}

func listHandler(cmd *cobra.Command, args []string) {
	// TODO Fetch details and display
	fmt.Println("filter : ", cmd.Flag("account").Value.String())
	fmt.Println("- service1")
	fmt.Println("- service2")
	fmt.Println("- service3")
}

func init() {
	List.Flags().StringP("account", "a", "", "Filter services based on the account address")
}
