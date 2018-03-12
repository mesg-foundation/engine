package cmdAccount

import (
	"fmt"

	"github.com/mesg-foundation/application/account"

	"github.com/spf13/cobra"
)

// List all the accounts
var List = &cobra.Command{
	Use:               "list",
	Short:             "List all accounts created on this computer",
	Example:           "mesg-cli account list",
	Run:               listHandler,
	DisableAutoGenTag: true,
}

func listHandler(cmd *cobra.Command, args []string) {
	for _, account := range account.List() {
		fmt.Println(account)
	}
}
