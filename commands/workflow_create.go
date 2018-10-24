package commands

import (
	"fmt"

	"github.com/mesg-foundation/core/utils/pretty"
	"github.com/spf13/cobra"
)

type createWorkflowCmd struct {
	baseCmd

	// name is an optionally set unique name for workflow.
	name string

	e WorkflowExecutor
}

func newCreateWorkflowCmd(e WorkflowExecutor) *createWorkflowCmd {
	c := &createWorkflowCmd{e: e}
	c.cmd = newCommand(&cobra.Command{
		Use:     "create WORKFLOW.yml",
		Short:   "Create and run a new workflow",
		Long:    "Create a new workflow from yaml file and run it",
		Example: `mesg-core workflow create WORKFLOW.yml`,
		Args:    cobra.ExactArgs(1),
		RunE:    c.runE,
	})
	c.cmd.Flags().StringVar(&c.name, "name", "", "unique name of workflow")
	return c
}

func (c *createWorkflowCmd) runE(cmd *cobra.Command, args []string) error {
	var (
		filePath   = args[0]
		workflowID string
		err        error
	)

	pretty.Progress("Creating workflow...", func() {
		workflowID, err = c.e.CreateWorkflow(filePath, c.name)
	})
	if err != nil {
		return err
	}

	fmt.Printf("%s Workflow running\n", pretty.SuccessSign)
	fmt.Printf("To see its logs, run the command:\n\tmesg-core workflow logs %s\n", workflowID)
	return nil
}
