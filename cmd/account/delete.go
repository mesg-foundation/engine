package cmdAccount

import (
	"fmt"

	"github.com/mesg-foundation/application/account"
	"github.com/mesg-foundation/application/cmd/utils"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// Delete a specific accounts
var Delete = &cobra.Command{
	Use:   "delete",
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
		account = detectAccount(args[0])
	}
	if account == nil {
		account = askAccount()
	}
	if cmdUtils.Confirm(cmd, "The account "+account.Name+" will be deleted. Are you sure ?") {
		// TODO add real deletion
		fmt.Println("account deleted", account)
	}
}

func detectAccount(accountValue string) *account.Account {
	accounts := account.List()
	for _, account := range accounts {
		if account.Name == accountValue || account.Address == accountValue || account.String() == accountValue {
			return account
		}
	}
	return nil
}

func askAccount() *account.Account {
	accounts := account.List()
	var selectedAccount string
	survey.AskOne(&survey.Select{
		Message: "Choose the account you want to delete",
		Options: accountOptions(accounts),
	}, &selectedAccount, nil)
	return detectAccount(selectedAccount)
}

func accountOptions(accounts []*account.Account) []string {
	res := make([]string, len(accounts))
	for i := 0; i < len(accounts); i++ {
		res[i] = accounts[i].String()
	}
	return res
}

func init() {
	cmdUtils.Confirmable(Delete)
}
