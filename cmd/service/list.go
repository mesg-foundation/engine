package cmdService

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	dbServices "github.com/mesg-foundation/core/database/services"
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
	services, err := dbServices.All() // TODO: this should use the API
	handleError(err)
	if len(services) == 0 {
		fmt.Println("No services")
		return
	}
	for _, service := range services {
		fmt.Println("-", aurora.Bold(service.Hash()), "-", service.Name)
	}
}
