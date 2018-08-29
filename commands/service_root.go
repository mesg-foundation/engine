package commands

import (
	"github.com/spf13/cobra"
)

type rootServiceCmd struct {
	baseCmd
}

func newRootServiceCmd(e ServiceExecutor) *rootServiceCmd {
	c := &rootServiceCmd{}
	c.cmd = newCommand(&cobra.Command{
		Use:   "service",
		Short: "Manage your services",
	})

	c.cmd.AddCommand(
		newServiceTestCmd().cmd,
		newServiceDeployCmd(e).cmd,
		newServiceValidateCmd(e).cmd,
		newServiceStartCmd(e).cmd,
		newServiceStopCmd(e).cmd,
		newServiceDetailCmd(e).cmd,
		newServiceListCmd(e).cmd,
		newServiceInitCmd(e).cmd,
		newServiceDeleteCmd(e).cmd,
		newServiceLogsCmd(e).cmd,
		newServiceDocsCmd(e).cmd,
		newServiceDevCmd(e).cmd,
		newServiceExecuteCmd(e).cmd,
	)
	return c
}
