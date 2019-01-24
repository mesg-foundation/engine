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
		Use:     "stop SERVICE",
		Short:   "Stop a service",
		Example: `mesg-core service stop SERVICE`,
		RunE:    c.runE,
	})
	return c
}

func (c *serviceStopCmd) runE(cmd *cobra.Command, args []string) error {
	for _, serviceID := range args {
		var err error
		pretty.Progress("Stopping service...", func() {
			err = c.e.ServiceStop(serviceID)
		})
		if err != nil {
			return err
		}
		fmt.Printf("%s Service %q stopped\n", pretty.SuccessSign, serviceID)
	}
	return nil
}
