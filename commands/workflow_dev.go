package commands

import (
	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/mesg-foundation/core/x/xsignal"
	"github.com/spf13/cobra"
)

type devWorkflowCmd struct {
	baseCmd

	e WorkflowExecutor
}

func newDevWorkflowCmd(e WorkflowExecutor) *devWorkflowCmd {
	c := &devWorkflowCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "dev WORKFLOW.yml",
		Short:   "Create and run a new workflow in development mode",
		Long:    "Create and run a new workflow in development mode which will be deleted after exiting",
		Example: `mesg-core workflow dev WORKFLOW.yml`,
		Args:    cobra.ExactArgs(1),
		RunE:    c.runE,
	})
	return c
}

func (c *devWorkflowCmd) runE(cmd *cobra.Command, args []string) error {
	var (
		filePath   = args[0]
		workflowID string
		err        error
	)

	pretty.Progress("Creating workflow...", func() {
		workflowID, err = c.e.CreateWorkflow(filePath, "")
	})
	if err != nil {
		return err
	}

	waitC, closer, err := newWorkflowLogsPrinter().Print(c.e, workflowID)
	if err != nil {
		return err
	}
	defer closer()

	select {
	case <-xsignal.WaitForInterrupt():
		closer()
		pretty.Progress("Removing the workflow...", func() {
			err = c.e.DeleteWorkflow(workflowID)
		})
		return err

	case err := <-waitC:
		return err
	}
}
