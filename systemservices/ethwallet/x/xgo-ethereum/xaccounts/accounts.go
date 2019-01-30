package xaccounts

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
)

// GetAccount retrieves an account from its address.
func GetAccount(keystore *keystore.KeyStore, address common.Address) (accounts.Account, error) {
	for _, account := range keystore.Accounts() {
		if account.Address == address {
			return account, nil
		}
	}
	return accounts.Account{}, fmt.Errorf("Account not found")
}
