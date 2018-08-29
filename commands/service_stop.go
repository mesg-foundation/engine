package commands

import (
	"fmt"

	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
)

type serviceStopCmd struct {
	baseCmd

	e ServiceExecutor
}

func newServiceStopCmd(e ServiceExecutor) *serviceStopCmd {
	c := &serviceStopCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:   "stop SERVICE",
		Short: "Stop a service",
		Long: `Stop a service.

**WARNING:** If you stop a service with your stake duration still ongoing, you may lost your stake.
You will **NOT** get your stake back immediately. You will get your remaining stake only after a delay.
To have more explanation, see the page [stake explanation from the documentation]().`, // TODO: add link
		Example: `mesg-core service stop SERVICE`,
		Args:    cobra.ExactArgs(1),
		RunE:    c.runE,
	})

	return c
}

func (c *serviceStopCmd) runE(cmd *cobra.Command, args []string) error {
	if err := c.e.ServiceStop(args[0]); err != nil {
		return err
	}
	fmt.Println(pretty.Success("Service stopped"))
	return nil
}
