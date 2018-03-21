package cmd

import (
	"github.com/mesg-foundation/application/cmd/marketplace/workflow"
	"github.com/spf13/cobra"
)

// MarketPlace manage all things from the marketplace
var MarketPlace = &cobra.Command{
	Use:               "marketplace",
	Short:             "Manage all things from the marketplace",
	DisableAutoGenTag: true,
}

// MarketPlaceWorkflow from the marketplace
var MarketPlaceWorkflow = &cobra.Command{
	Use:               "workflow",
	Short:             "Manage workflows from the marketplace",
	DisableAutoGenTag: true,
}

func init() {
	MarketPlaceWorkflow.AddCommand(cmdWorkflowMarketPlace.Detail)
	MarketPlaceWorkflow.AddCommand(cmdWorkflowMarketPlace.List)

	MarketPlace.AddCommand(MarketPlaceWorkflow)

	RootCmd.AddCommand(MarketPlace)
}
