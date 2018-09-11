package commands

import (
	"fmt"
	"os"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
)

type logsCmd struct {
	baseCmd

	e RootExecutor
}

func newLogsCmd(e RootExecutor) *logsCmd {
	c := &logsCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:   "logs",
		Short: "Show the Core's logs",
		RunE:  c.runE,
	})
	return c
}

func (c *logsCmd) runE(cmd *cobra.Command, args []string) error {
	status, err := c.e.Status()
	if err != nil {
		return err
	}

	if status == container.STOPPED {
		fmt.Println(pretty.Warn("Core is stopped"))
		return nil
	}

	reader, err := c.e.Logs()
	if err != nil {
		return err
	}
	defer reader.Close()

	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, reader)
	return err
}
