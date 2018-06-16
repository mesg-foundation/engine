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
	Service.AddCommand(service.Deploy)
	Service.AddCommand(service.Validate)
	Service.AddCommand(service.Test)
	Service.AddCommand(service.Start)
	Service.AddCommand(service.Stop)
	// Service.AddCommand(service.Pause)
	// Service.AddCommand(service.Resume)
	Service.AddCommand(service.Detail)
	Service.AddCommand(service.List)
	Service.AddCommand(service.Status)
	Service.AddCommand(service.Init)
	Service.AddCommand(service.Delete)

	RootCmd.AddCommand(Service)
}
