package account

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/application/config"
)

// TODO add real account creation
func generate(password string, name string) (address common.Address, err error) {
	if password == "" {
		err = errors.New("Password is missing")
		return
	}
	if name == "" {
		err = errors.New("Name is missing")
		return
	}

	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP

	address, err = keystore.StoreKey(config.AccountDirectory, password, scryptN, scryptP)

	return
}

// Generate an account based on some predefined data
func (account *Account) Generate() (err error) {
	addr, err := generate(account.Password, account.Name)
	if err != nil {
		return
	}
	account.Address = addr
	return
}
