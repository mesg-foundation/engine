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
		Args:    cobra.ExactArgs(1),
		RunE:    c.runE,
	})
	return c
}

func (c *serviceStopCmd) runE(cmd *cobra.Command, args []string) error {
	var err error
	pretty.Progress("Stopping service...", func() {
		err = c.e.ServiceStop(args[0])
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s Service stopped\n", pretty.SuccessSign)
	return nil
}
