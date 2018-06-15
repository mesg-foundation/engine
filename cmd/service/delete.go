package cmdService

import (
	"context"
	"fmt"

	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/spf13/cobra"
)

// Delete a service to the marketplace
var Delete = &cobra.Command{
	Use:               "delete",
	Short:             "Delete a service",
	Example:           `mesg-core service delete`,
	Run:               deleteHandler,
	DisableAutoGenTag: true,
}

func deleteHandler(cmd *cobra.Command, args []string) {
	for _, arg := range args {
		_, err := cli.DeleteService(context.Background(), &core.DeleteServiceRequest{
			ServiceID: arg,
		})
		cmdUtils.HandleError(err)
		fmt.Println("Service", arg, "deleted")
	}
}
