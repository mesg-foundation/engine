package cmdUtils

import (
	"github.com/mesg-foundation/application/account"
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

func accountOptions(accounts []*account.Account) []string {
	res := make([]string, len(accounts))
	for i := 0; i < len(accounts); i++ {
		res[i] = accounts[i].String()
	}
	return res
}
