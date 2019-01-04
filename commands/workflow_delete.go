package commands

import (
	"fmt"

	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
)

type deleteWorkflowCmd struct {
	baseCmd

	e WorkflowExecutor
}

func newDeleteWorkflowCmd(e WorkflowExecutor) *deleteWorkflowCmd {
	c := &deleteWorkflowCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "delete WORKFLOW",
		Short:   "Stop and delete a workflow",
		Example: `mesg-core workflow delete WORKFLOW`,
		Args:    cobra.ExactArgs(1),
		RunE:    c.runE,
	})
	return c
}

func (c *deleteWorkflowCmd) runE(cmd *cobra.Command, args []string) error {
	var (
		workflowID = args[0]
		err        error
	)

	pretty.Progress(fmt.Sprintf("Deleting workflow %q...", workflowID), func() {
		err = c.e.DeleteWorkflow(workflowID)
	})
	if err != nil {
		return err
	}

	fmt.Printf("%s workflow %q deleted\n", pretty.SuccessSign, workflowID)
	return nil
}
