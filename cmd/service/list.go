package cmdService

import (
	"fmt"

	"github.com/mesg-foundation/core/database/services"
	"github.com/spf13/cobra"
)

// List all the services
var List = &cobra.Command{
	Use:   "list",
	Short: "List all published services",
	Long: `This command returns all published services with basic information.
Optionally, you can filter the services published by a specific developer:
To have more details, see the [detail command](mesg-core_service_detail.md).`,
	Example:           `mesg-core service list`,
	Run:               listHandler,
	DisableAutoGenTag: true,
}

func listHandler(cmd *cobra.Command, args []string) {
	services, err := services.All()
	handleError(err)
	for _, service := range services {
		fmt.Println("-", service.Name)
	}
}
