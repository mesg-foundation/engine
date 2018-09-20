package commands

import (
	"fmt"
	"strings"

	"github.com/mesg-foundation/core/protobuf/coreapi"

	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
)

type serviceDetailCmd struct {
	baseCmd
	e ServiceExecutor
}

func newServiceDetailCmd(e ServiceExecutor) *serviceDetailCmd {
	c := &serviceDetailCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "detail SERVICE",
		Short:   "Show details of a published service",
		Args:    cobra.ExactArgs(1),
		Example: "mesg-core service detail SERVICE",
		RunE:    c.runE,
	})
	return c
}

func (c *serviceDetailCmd) runE(cmd *cobra.Command, args []string) error {
	var (
		err     error
		service *coreapi.Service
	)
	pretty.Progress("Loading service...", func() {
		service, err = c.e.ServiceByID(args[0])
	})
	if err != nil {
		return err
	}
	fmt.Printf("name: %s\n", pretty.Bold(service.Name))
	fmt.Println("events:")
	for _, event := range service.Events {
		params := []string{}
		for _, d := range event.Data {
			params = append(params, d.Key+" "+d.Type)
		}
		fmt.Printf("  %s (%s)\n", pretty.Bold(event.Key), strings.Join(params, ", "))
	}
	fmt.Println("tasks:")
	for _, task := range service.Tasks {
		fmt.Printf("  %s:\n", pretty.Bold(task.Key))
		for _, output := range task.Outputs {
			params := []string{}
			for _, param := range output.Data {
				params = append(params, param.Key+" "+param.Type)
			}
			fmt.Printf("    %s (%s)\n", pretty.Bold(output.Key), strings.Join(params, ", "))
		}
	}
	return nil
}
