package cmdUtils

import (
	accountPkg "github.com/mesg-foundation/application/account"
	"github.com/spf13/cobra"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// FindAccount returns an account if it matches either the address or the name
func FindAccount(accountValue string) (account *accountPkg.Account) {
	accounts := accountPkg.List()
	for _, a := range accounts {
		if a.Name == accountValue || a.Address == accountValue || a.String() == accountValue {
			account = a
			return
		}
	}
	return
}

// AskAccount asks to the user to select an account
func AskAccount(message string) (account *accountPkg.Account) {
	accounts := accountPkg.List()
	var selectedAccount string
	survey.AskOne(&survey.Select{
		Message: message,
		Options: accountOptions(accounts),
	}, &selectedAccount, nil)
	account = FindAccount(selectedAccount)
	return
}

// AccountFromFlag returns the account based on the flag
func AccountFromFlag(cmd *cobra.Command) (account *accountPkg.Account) {
	account = FindAccount(cmd.Flag("account").Value.String())
	return
}

// AccountFromFlagOrAsk return the selected account either by reading the flag value or by asking to the user
func AccountFromFlagOrAsk(cmd *cobra.Command, message string) (account *accountPkg.Account) {
	account = AccountFromFlag(cmd)
	if account == nil {
		account = AskAccount(message)
	}
	return
}

func accountOptions(accounts []*accountPkg.Account) (options []string) {
	options = make([]string, len(accounts))
	for i := 0; i < len(accounts); i++ {
		options[i] = accounts[i].String()
	}
	return
}

// Accountable mark a command as Accountable so will have the --account flag
func Accountable(cmd *cobra.Command) {
	cmd.Flags().StringP("account", "a", "", "Account you want to use")
}
