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
		Example: `mesg-core service stop SERVICE [SERVICE...]`,
		RunE:    c.runE,
	})
	return c
}

func (c *serviceStopCmd) runE(cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		var err error
		// build function to avoid using arg inside progress
		fn := func(serviceID string) func() {
			return func() {
				err = c.e.ServiceStop(serviceID)
			}
		}(arg)
		pretty.Progress(fmt.Sprintf("Stopping service %q...", arg), fn)
		if err != nil {
			return err
		}
		fmt.Printf("%s Service %q stopped\n", pretty.SuccessSign, arg)
	}
	return nil
}
