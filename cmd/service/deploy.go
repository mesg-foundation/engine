package service

import (
	"context"
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/spf13/cobra"
)

// Deploy a service to the marketplace.
var Deploy = &cobra.Command{
	Use:     "deploy",
	Aliases: []string{"publish"},
	Short:   "Deploy a service",
	Long: `Deploy a service.

To get more information, see the [deploy page from the documentation](https://docs.mesg.com/guide/service/deploy-a-service.html)`,
	Example:           `mesg-core service deploy PATH_TO_SERVICE`,
	Run:               deployHandler,
	DisableAutoGenTag: true,
}

func deployHandler(cmd *cobra.Command, args []string) {
	service := prepareService(defaultPath(args))
	reply, err := cli().DeployService(context.Background(), &core.DeployServiceRequest{
		Service: service,
	})
	utils.HandleError(err)
	fmt.Println("Service deployed with ID:", aurora.Green(reply.ServiceID))
	fmt.Println("To start it, run the command: mesg-core service start " + reply.ServiceID)
}
