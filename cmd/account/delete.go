package cmdAccount

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"
	"github.com/spf13/cobra"
)

// Delete a specific accounts
var Delete = &cobra.Command{
	Use:   "delete",
	Short: "Delete an account",
	Example: `mesg-cli service delete --account 0x0000000000000000000000000000000000000000
mesg-cli service delete`,
	Run:               deleteHandler,
	DisableAutoGenTag: true,
}

func deleteHandler(cmd *cobra.Command, args []string) {
	acc := cmdUtils.AccountFromFlagOrAsk(cmd, "Choose the account you want to delete")
	if cmdUtils.Confirm(cmd, "The account "+acc.Address.String()+" will be deleted. Are you sure ?") {
		// TODO add real deletion
		fmt.Println("account deleted", acc)
	}
}

func init() {
	cmdUtils.Confirmable(Delete)
	cmdUtils.Accountable(Delete)
}
