package account

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/mesg-foundation/application/config"
)

// List all available accounts on this computer
func List() (accountList []*Account) {
	store := keystore.NewKeyStore(config.AccountDirectory, keystore.StandardScryptN, keystore.StandardScryptP)
	for _, wallet := range store.Wallets() {
		for _, account := range wallet.Accounts() {
			accountList = append(accountList, &Account{
				Address: account.Address,
				URL:     account.URL.String(),
			})
		}
	}
	return accountList
}
