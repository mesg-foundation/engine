package cmd

import (
	"github.com/mesg-foundation/application/cmd/account"
	"github.com/spf13/cobra"
)

// Account is the command to manage all account activity
var Account = &cobra.Command{
	Use:               "account",
	Short:             "Manage your MESG accounts",
	DisableAutoGenTag: true,
}

func init() {
	Account.AddCommand(cmdAccount.Create)
	Account.AddCommand(cmdAccount.List)
	Account.AddCommand(cmdAccount.Delete)

	RootCmd.AddCommand(Account)
}
