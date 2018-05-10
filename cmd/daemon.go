package cmd

import (
	"github.com/mesg-foundation/core/cmd/daemon"
	"github.com/spf13/cobra"
)

// Daemon is the root command related to the daemon
var Daemon = &cobra.Command{
	Use:               "daemon",
	Short:             "Manage the MESG daemon",
	DisableAutoGenTag: true,
}

func init() {
	Daemon.AddCommand(daemon.Start)
	Daemon.AddCommand(daemon.Stop)
	Daemon.AddCommand(daemon.Status)
	Daemon.AddCommand(daemon.Logs)

	RootCmd.AddCommand(Daemon)
}
