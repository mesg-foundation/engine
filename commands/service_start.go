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
		Long:    "Start a service previously published services.",
		Example: `mesg-core service start SERVICE`,
		Args:    cobra.ExactArgs(1),
		RunE:    c.runE,
	})
	return c
}

func (c *serviceStartCmd) runE(cmd *cobra.Command, args []string) error {
	var (
		serviceID = args[0]
		err       error
	)
	pretty.Progress("Starting service...", func() {
		err = c.e.ServiceStart(serviceID)
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s Service started.\n", pretty.SuccessSign)
	fmt.Printf("To see its logs, run the command:\n\tmesg-core service logs %s\n", serviceID)
	return nil
}
