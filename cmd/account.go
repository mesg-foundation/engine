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
	// Account.AddCommand(account.Create)
	// Account.AddCommand(account.List)
	// Account.AddCommand(account.Detail)
	// Account.AddCommand(account.Delete)
	// Account.AddCommand(account.Export)
	// Account.AddCommand(account.Import)

	// RootCmd.AddCommand(Account)
}
