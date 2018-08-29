package commands

import (
	"fmt"
	"strings"

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
	service, err := c.e.ServiceByID(args[0])
	if err != nil {
		return err
	}
	fmt.Println("name: ", pretty.Bold(service.Name))
	fmt.Println("events: ")
	for name, event := range service.Events {
		params := []string{}
		for key, d := range event.Data {
			params = append(params, key+" "+d.Type)
		}
		fmt.Println("  ", pretty.Bold(name), "(", strings.Join(params, ", "), ")")
	}
	fmt.Println("tasks: ")
	for name, task := range service.Tasks {
		fmt.Println("  ", pretty.Bold(name), ":")
		for outputName, output := range task.Outputs {
			params := []string{}
			for paramName, param := range output.Data {
				params = append(params, paramName+" "+param.Type)
			}
			fmt.Println("    ", pretty.Bold(outputName), "(", strings.Join(params, ", "), ")")
		}
	}
	return nil
}
