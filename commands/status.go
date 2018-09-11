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
		Short: "Status of the Core",
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
		fmt.Println(pretty.Success("Core is running"))
	} else {
		fmt.Println(pretty.Warn("Core is stopped"))
	}
	return nil
}
