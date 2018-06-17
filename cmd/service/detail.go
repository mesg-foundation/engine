package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/cmd/utils"
	"github.com/spf13/cobra"
)

// Detail returns all the details of a service
var Detail = &cobra.Command{
	Use:               "detail SERVICE_ID",
	Short:             "Show details of a published service",
	Args:              cobra.MinimumNArgs(1),
	Example:           "mesg-core service detail SERVICE_ID",
	Run:               detailHandler,
	DisableAutoGenTag: true,
}

func detailHandler(cmd *cobra.Command, args []string) {
	serviceReply, err := cli.Service(context.Background(), &core.ServiceRequest{
		ServiceID: args[0],
	})
	utils.HandleError(err)
	service := serviceReply.Service
	fmt.Println("name: ", aurora.Bold(service.Name))
	fmt.Println("events: ")
	for name, event := range service.Events {
		params := []string{}
		for key, d := range event.Data {
			params = append(params, key+" "+d.Type)
		}
		fmt.Println("  ", aurora.Bold(name), "(", strings.Join(params, ", "), ")")
	}
	fmt.Println("tasks: ")
	for name, task := range service.Tasks {
		fmt.Println("  ", aurora.Bold(name), ":")
		for outputName, output := range task.Outputs {
			params := []string{}
			for paramName, param := range output.Data {
				params = append(params, paramName+" "+param.Type)
			}
			fmt.Println("    ", aurora.Bold(outputName), "(", strings.Join(params, ", "), ")")
		}
	}
}
