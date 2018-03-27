package cmdService

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Detail returns all the details of a service
var Detail = &cobra.Command{
	Use:               "detail SERVICE_ID",
	Short:             "Show details of a published service",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli service detail SERVICE_ID",
	Run:               detailHandler,
	DisableAutoGenTag: true,
}

func detailHandler(cmd *cobra.Command, args []string) {
	// TODO Fetch details and display
	fmt.Println("service details : ", args)
}
