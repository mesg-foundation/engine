package cmdService

import (
	"context"
	"fmt"

	"github.com/mesg-foundation/core/api/core"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

// Publish a service to the marketplace
var Publish = &cobra.Command{
	Use:   "publish",
	Short: "Publish a service",
	Long: `Publish a service on the Network.

To get more information, see the [publish page from the documentation](https://docs.mesg.tech/service/develop/publish.html)`,
	Example:           `mesg-core service publish`,
	Run:               deployHandler,
	DisableAutoGenTag: true,
}

func deployHandler(cmd *cobra.Command, args []string) {
	service := loadService(defaultPath(args))

	buildDockerImage(defaultPath(args), service.Name)

	reply, err := cli.DeployService(context.Background(), &core.DeployServiceRequest{
		Service: service,
	})
	handleError(err)
	fmt.Println("Service deployed with ID:", aurora.Green(reply.ServiceID))
}
