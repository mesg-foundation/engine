package cmd

import (
	"github.com/mesg-foundation/core/cmd/service"
	"github.com/spf13/cobra"
)

// Service is the root command related to services
var Service = &cobra.Command{
	Use:               "service",
	Short:             "Manage your services",
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
	Service.AddCommand(cmdService.Status)
	Service.AddCommand(cmdService.Init)

	RootCmd.AddCommand(Service)
}
