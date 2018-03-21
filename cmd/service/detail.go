package cmdService

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Detail returns all the details of a service
var Detail = &cobra.Command{
	Use:               "detail SERVICE",
	Short:             "Show details of a service",
	Long:              "Provide details about a service like number of nodes running it, average revenue, etc..",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli marketplace service detail ethereum",
	Run:               detailHandler,
	DisableAutoGenTag: true,
}

func detailHandler(cmd *cobra.Command, args []string) {
	// TODO Fetch details and display
	fmt.Println("service details : ", args)
}
