package cmd

import (
	"github.com/mesg-foundation/application/cmd/service"

	"github.com/spf13/cobra"
)

// Service is the root command related to services
var Service = &cobra.Command{
	Use:               "service",
	Short:             "Manage the services you are running",
	DisableAutoGenTag: true,
}

func init() {
	Service.AddCommand(cmdService.Deploy)

	RootCmd.AddCommand(Service)
}
