package commands

import (
	"fmt"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
)

type statusCmd struct {
	baseCmd

	e RootExecutor
}

func newStatusCmd(e RootExecutor) *statusCmd {
	c := &statusCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:   "status",
		Short: "Get the Core's status",
		RunE:  c.runE,
	})
	return c
}

func (c *statusCmd) runE(cmd *cobra.Command, args []string) error {
	status, err := c.e.Status()
	if err != nil {
		return err
	}

	if status == container.RUNNING {
		fmt.Printf("%s Core is running\n", pretty.SuccessSign)
	} else {
		fmt.Printf("%s Core is stopped\n", pretty.WarnSign)
	}
	return nil
}
