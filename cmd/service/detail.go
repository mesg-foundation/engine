package cmdService

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Detail returns all the details of a service
var Detail = &cobra.Command{
	Use:               "detail SERVICE_FOLDER",
	Short:             "Show details of a published service",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-cli service detail SERVICE_FOLDER",
	Run:               detailHandler,
	DisableAutoGenTag: true,
}

func detailHandler(cmd *cobra.Command, args []string) {
	service := loadService(defaultPath(args))
	fmt.Println("service details : ", service)
}
