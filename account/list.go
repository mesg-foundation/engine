package account

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/mesg-foundation/core/config"
)

// List all available accounts on this computer
func List() (accountList []accounts.Account) {
	for _, wallet := range config.Store.Wallets() {
		for _, a := range wallet.Accounts() {
			accountList = append(accountList, a)
		}
	}
	return accountList
}
