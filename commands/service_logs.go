package commands

import (
	"os"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/mesg-foundation/core/x/xsignal"
	"github.com/spf13/cobra"
)

type serviceLogsCmd struct {
	baseCmd

	dependency string

	e ServiceExecutor
}

func newServiceLogsCmd(e ServiceExecutor) *serviceLogsCmd {
	c := &serviceLogsCmd{
		e:          e,
		dependency: "*",
	}

	c.cmd = newCommand(&cobra.Command{
		Use:   "logs",
		Short: "Show the logs of a service",
		Example: `mesg-core service logs SERVICE
mesg-core service logs SERVICE --dependency DEPENDENCY_NAME`,
		Args: cobra.ExactArgs(1),
		RunE: c.runE,
	})
	c.cmd.Flags().StringVarP(&c.dependency, "dependency", "d", c.dependency, "Name of the dependency to show the logs from")
	return c
}

func (c *serviceLogsCmd) runE(cmd *cobra.Command, args []string) error {
	readers, err := c.e.ServiceDependencyLogs(args[0], c.dependency)
	if err != nil {
		return err
	}
	defer func() {
		for _, reader := range readers {
			reader.Close()
		}
	}()

	for _, reader := range readers {
		go stdcopy.StdCopy(os.Stdout, os.Stderr, reader)
	}

	<-xsignal.WaitForInterrupt()
	return nil
}
