package cmdAccount

import (
	"fmt"

	"github.com/mesg-foundation/application/account"
	"github.com/mesg-foundation/application/cmd/utils"
	"github.com/spf13/cobra"
)

// Delete a specific accounts
var Delete = &cobra.Command{
	Use:   "delete ACCOUNT",
	Short: "Delete an account",
	Example: `mesg-cli service delete accountX
mesg-cli service delete 0x0000000000000000000000000000000000000000
mesg-cli service delete`,
	Run:               deleteHandler,
	DisableAutoGenTag: true,
}

func deleteHandler(cmd *cobra.Command, args []string) {
	var account *account.Account
	if len(args) > 0 {
		account = cmdUtils.FindAccount(args[0])
	}
	if account == nil {
		account = cmdUtils.AskAccount("Choose the account you want to delete")
	}
	if cmdUtils.Confirm(cmd, "The account "+account.Name+" will be deleted. Are you sure ?") {
		// TODO add real deletion
		fmt.Println("account deleted", account)
	}
}

func init() {
	cmdUtils.Confirmable(Delete)
}
