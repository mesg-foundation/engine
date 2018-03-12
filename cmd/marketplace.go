package cmd

import (
	"github.com/mesg-foundation/application/cmd/marketplace/service"
	"github.com/mesg-foundation/application/cmd/marketplace/workflow"
	"github.com/spf13/cobra"
)

// MarketPlace manage all things from the marketplace
var MarketPlace = &cobra.Command{
	Use:               "marketplace",
	Short:             "Manage all things from the marketplace",
	DisableAutoGenTag: true,
}

// MarketPlaceService from the marketplace
var MarketPlaceService = &cobra.Command{
	Use:               "service",
	Short:             "Manage services from the marketplace",
	DisableAutoGenTag: true,
}

// MarketPlaceWorkflow from the marketplace
var MarketPlaceWorkflow = &cobra.Command{
	Use:               "workflow",
	Short:             "Manage workflows from the marketplace",
	DisableAutoGenTag: true,
}

func init() {
	MarketPlaceService.AddCommand(cmdServiceMarketPlace.Start)
	MarketPlaceService.AddCommand(cmdServiceMarketPlace.Stop)
	MarketPlaceService.AddCommand(cmdServiceMarketPlace.Pause)
	MarketPlaceService.AddCommand(cmdServiceMarketPlace.Resume)

	MarketPlaceWorkflow.AddCommand(cmdWorkflowMarketPlace.Detail)
	MarketPlaceWorkflow.AddCommand(cmdWorkflowMarketPlace.List)

	MarketPlace.AddCommand(MarketPlaceService)
	MarketPlace.AddCommand(MarketPlaceWorkflow)

	RootCmd.AddCommand(MarketPlace)
}
