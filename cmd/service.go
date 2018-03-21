package cmd

import (
	"github.com/mesg-foundation/application/cmd/service"
	"github.com/spf13/cobra"
)

// Service is the root command related to services
var Service = &cobra.Command{
	Use:               "service",
	Short:             "Manage the services you created",
	DisableAutoGenTag: true,
}

func init() {
	Service.AddCommand(cmdService.Publish)
	Service.AddCommand(cmdService.Validate)
	Service.AddCommand(cmdService.Test)
	Service.AddCommand(cmdService.Start)
	Service.AddCommand(cmdService.Stop)
	Service.AddCommand(cmdService.Pause)
	Service.AddCommand(cmdService.Resume)
	Service.AddCommand(cmdService.Detail)
	Service.AddCommand(cmdService.List)

	RootCmd.AddCommand(Service)
}
