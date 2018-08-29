package commands

import (
	"fmt"

	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
)

type serviceStartCmd struct {
	baseCmd

	e ServiceExecutor
}

func newServiceStartCmd(e ServiceExecutor) *serviceStartCmd {
	c := &serviceStartCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "start SERVICE",
		Short:   "Start a service",
		Long:    "Start a service from the published available services. You have to provide a stake value and duration.",
		Example: `mesg-core service start SERVICE`,
		Args:    cobra.ExactArgs(1),
		RunE:    c.runE,
	})
	return c
}

func (c *serviceStartCmd) runE(cmd *cobra.Command, args []string) error {
	if err := c.e.ServiceStart(args[0]); err != nil {
		return err
	}

	fmt.Println(pretty.Success("Service is running"))
	return nil
}
