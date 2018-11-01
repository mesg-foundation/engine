package commands

import (
	"github.com/spf13/cobra"
)

type rootWorkflowCmd struct {
	baseCmd
}

func newRootWorkflowCmd(e WorkflowExecutor) *rootWorkflowCmd {
	c := &rootWorkflowCmd{}
	c.cmd = newCommand(&cobra.Command{
		Use:   "workflow",
		Short: "Manage workflows",
	})

	c.cmd.AddCommand(
		newCreateWorkflowCmd(e).cmd,
		newDeleteWorkflowCmd(e).cmd,
		newWorkflowLogsCmd(e).cmd,
		newDevWorkflowCmd(e).cmd,
	)

	return c
}
