package cmd

import (
	"github.com/mesg-foundation/application/cmd/account"
	"github.com/spf13/cobra"
)

// Account are a group of commands to manage your accounts
var Account = &cobra.Command{
	Use:               "account",
	Short:             "Manage your accounts",
	DisableAutoGenTag: true,
}

func init() {
	Account.AddCommand(cmdAccount.Create)
	Account.AddCommand(cmdAccount.List)
	Account.AddCommand(cmdAccount.Detail)
	Account.AddCommand(cmdAccount.Delete)
	Account.AddCommand(cmdAccount.Export)
	Account.AddCommand(cmdAccount.Import)

	RootCmd.AddCommand(Account)
}
