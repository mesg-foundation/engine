package cmdAccount

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/application/account"
	"github.com/mesg-foundation/application/cmd/utils"
	"github.com/spf13/cobra"
)

// Delete a specific accounts
var Delete = &cobra.Command{
	Use:   "delete",
	Short: "Delete an account",
	Long: `This method deletes an account.

**Warning:** If you didn't previously [export this account](mesg-cli_account_export.md), you will lost it **forever.**`,
	Example: `mesg-cli service delete
mesg-cli service delete --account ACCOUNT --confirm`,
	Run:               deleteHandler,
	DisableAutoGenTag: true,
}

func deleteHandler(cmd *cobra.Command, args []string) {
	acc := cmdUtils.AccountFromFlagOrAsk(cmd, "Choose the account to delete:")
	if cmdUtils.Confirm(cmd, "The account "+acc.Address.String()+" will be deleted. Are you sure?") {
		if err := account.Destroy(acc); err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", aurora.Green("Account deleted with success").Bold())
	}
}

func init() {
	cmdUtils.Confirmable(Delete)
	cmdUtils.Accountable(Delete)
}
