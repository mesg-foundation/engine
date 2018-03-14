package cmdUtils

import (
	"github.com/mesg-foundation/application/account"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// FindAccount returns an account if it matches either the address or the name
func FindAccount(accountValue string) *account.Account {
	accounts := account.List()
	for _, account := range accounts {
		if account.Name == accountValue || account.Address == accountValue || account.String() == accountValue {
			return account
		}
	}
	return nil
}

// AskAccount asks to the user to select an account
func AskAccount(message string) *account.Account {
	accounts := account.List()
	var selectedAccount string
	survey.AskOne(&survey.Select{
		Message: message,
		Options: accountOptions(accounts),
	}, &selectedAccount, nil)
	return FindAccount(selectedAccount)
}

// AccountFromFlag returns the account based on the flag
func AccountFromFlag(cmd *cobra.Command) *account.Account {
	return FindAccount(cmd.Flag("account").Value.String())
}

// AccountFromFlagOrAsk return the selected account either by reading the flag value or by asking to the user
func AccountFromFlagOrAsk(cmd *cobra.Command, message string) *account.Account {
	account := AccountFromFlag(cmd)
	if account == nil {
		account = AskAccount(message)
	}
	return account
}

func accountOptions(accounts []*account.Account) []string {
	res := make([]string, len(accounts))
	for i := 0; i < len(accounts); i++ {
		res[i] = accounts[i].String()
	}
	return res
}

// Accountable mark a command as Accountable so will have the --account flag
func Accountable(cmd *cobra.Command) *cobra.Command {
	cmd.Flags().StringP("account", "a", "", "Account you want to use")
	return cmd
}
