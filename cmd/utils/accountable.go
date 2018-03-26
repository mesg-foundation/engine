package cmdUtils

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/mesg-foundation/application/account"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// FindAccount returns an account if it matches either the address or the name
func FindAccount(accountValue string) (acc accounts.Account) {
	accounts := account.List()
	for _, a := range accounts {
		if a.Address.String() == accountValue {
			acc = a
			break
		}
	}
	return
}

// AskAccount asks to the user to select an account
func AskAccount(message string) (acc accounts.Account) {
	accounts := account.List()
	if len(accounts) == 0 {
		panic(errors.New("Create an account first"))
	}
	var selectedAccount string
	survey.AskOne(&survey.Select{
		Message: message,
		Options: accountOptions(accounts),
	}, &selectedAccount, nil)
	acc = FindAccount(selectedAccount)
	return
}

// AccountFromFlag returns the account based on the flag
func AccountFromFlag(cmd *cobra.Command) (acc accounts.Account) {
	acc = FindAccount(cmd.Flag("account").Value.String())
	return
}

// AccountFromFlagOrAsk return the selected account either by reading the flag value or by asking to the user
func AccountFromFlagOrAsk(cmd *cobra.Command, message string) (acc accounts.Account) {
	acc = AccountFromFlag(cmd)
	if acc == (accounts.Account{}) {
		acc = AskAccount(message)
	}
	return
}

func accountOptions(accounts []accounts.Account) (options []string) {
	options = make([]string, len(accounts))
	for i := 0; i < len(accounts); i++ {
		options[i] = accounts[i].Address.String()
	}
	return
}

// Accountable mark a command as Accountable so will have the --account flag
func Accountable(cmd *cobra.Command) {
	cmd.Flags().StringP("account", "a", "", "Account you want to use")
}
