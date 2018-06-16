package cmd

import (
	"github.com/spf13/cobra"
)

// Workflow is the root command related to workflows
var Workflow = &cobra.Command{
	Use:               "workflow",
	Short:             "Manage your workflows",
	DisableAutoGenTag: true,
}

// TODO this command is disabled for now waiting for the Workflow feature to be implemented
func init() {

	// Workflow.AddCommand(workflow.Pause)
	// Workflow.AddCommand(workflow.Resume)
	// Workflow.AddCommand(workflow.Deploy)
	// Workflow.AddCommand(workflow.Test)
	// Workflow.AddCommand(workflow.Validate)
	// Workflow.AddCommand(workflow.List)
	// Workflow.AddCommand(workflow.Kill)
	// Workflow.AddCommand(workflow.Topup)
	// Workflow.AddCommand(workflow.Log)
	// Workflow.AddCommand(workflow.Detail)

	// RootCmd.AddCommand(Workflow)
}
