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
		Long:    "Start a service previously published services",
		Example: `mesg-core service start SERVICE`,
		RunE:    c.runE,
	})
	return c
}

func (c *serviceStartCmd) runE(cmd *cobra.Command, args []string) error {
	for _, serviceID := range args {
		var err error
		pretty.Progress("Starting service...", func() {
			err = c.e.ServiceStart(serviceID)
		})
		if err != nil {
			return err
		}
		fmt.Printf("%s Service %q started\n", pretty.SuccessSign, serviceID)
		if len(args) == 1 {
			fmt.Printf("To see its logs, run the command:\n\tmesg-core service logs %s\n", serviceID)
		}
	}
	return nil
}
