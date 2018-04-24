package cmdAccount

import (
	"fmt"

	"github.com/mesg-foundation/core/account"
	"github.com/spf13/cobra"
)

// List all the accounts
var List = &cobra.Command{
	Use:               "list",
	Short:             "List your local accounts",
	Example:           "mesg-cli account list",
	Run:               listHandler,
	DisableAutoGenTag: true,
}

func listHandler(cmd *cobra.Command, args []string) {
	for _, account := range account.List() {
		// TODO: can facto with a displaySummary function like in create.go
		fmt.Println(account.Address.String())
	}
}
