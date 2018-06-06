package cmdService

import (
	"context"
	"fmt"

	"github.com/mesg-foundation/core/api/core"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// Deploy a service to the marketplace
var Deploy = &cobra.Command{
	Use:     "deploy",
	Aliases: []string{"publish"},
	Short:   "Deploy a service",
	Long: `Deploy a service on the Network.

To get more information, see the [deploy page from the documentation](https://docs.mesg.tech/service/publish-a-service)`,
	Example:           `mesg-core sevice deploy`,
	Run:               deployHandler,
	DisableAutoGenTag: true,
}

func deployHandler(cmd *cobra.Command, args []string) {
	service := prepareService(defaultPath(args))
	reply, err := cli.DeployService(context.Background(), &core.DeployServiceRequest{
		Service: service,
	})
	handleError(err)
	fmt.Println("Service deployed with ID:", aurora.Green(reply.ServiceID))
}
