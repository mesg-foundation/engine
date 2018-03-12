package cmdAccount

import (
	"fmt"

	"github.com/mesg-foundation/application/cmd/utils"
	"github.com/mesg-foundation/application/types"
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
	var account *types.Account
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

func detectAccount(accountValue string) *types.Account {
	accounts := accountList()
	for _, account := range accounts {
		if account.Name == accountValue || account.Address == accountValue || account.String() == accountValue {
			return account
		}
	}
	return nil
}

func askAccount() *types.Account {
	accounts := accountList()
	var selectedAccount string
	survey.AskOne(&survey.Select{
		Message: "Choose the account you want to delete",
		Options: accountOptions(accounts),
	}, &selectedAccount, nil)
	return detectAccount(selectedAccount)
}

func accountList() []*types.Account {
	// TODO add real list
	return []*types.Account{
		&types.Account{Name: "Test1", Address: "0x0000000000000000000000000000000000000000"},
		&types.Account{Name: "Test2", Address: "0x0000000000000000000000000000000000000001"},
	}
}

func accountOptions(accounts []*types.Account) []string {
	res := make([]string, len(accounts))
	for i := 0; i < len(accounts); i++ {
		res[i] = accounts[i].String()
	}
	return res
}

func init() {
	cmdUtils.Confirmable(Delete)
}
