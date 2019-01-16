package xaccounts

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
)

// GetAccount retrieves an account from its address.
func GetAccount(keystore *keystore.KeyStore, address string) (account accounts.Account, err error) {
	err = fmt.Errorf("Account not found")
	_address := common.HexToAddress(address)
	for _, _account := range keystore.Accounts() {
		if _account.Address == _address {
			account = _account
			err = nil
			break
		}
	}
	return account, err
}
