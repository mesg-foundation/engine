package cmd

import (
	"github.com/spf13/cobra"
)

// Account are a group of commands to manage your accounts
var Account = &cobra.Command{
	Use:               "account",
	Short:             "Manage your accounts",
	DisableAutoGenTag: true,
}

// TODO this command is disabled for now waiting to have needs for account
func init() {
	// Account.AddCommand(cmdAccount.Create)
	// Account.AddCommand(cmdAccount.List)
	// Account.AddCommand(cmdAccount.Detail)
	// Account.AddCommand(cmdAccount.Delete)
	// Account.AddCommand(cmdAccount.Export)
	// Account.AddCommand(cmdAccount.Import)

	// RootCmd.AddCommand(Account)
}
