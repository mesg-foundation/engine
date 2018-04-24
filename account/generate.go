package account

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/mesg-foundation/core/config"
)

// Generate create a new account based on a password
func Generate(password string) (acc accounts.Account, err error) {
	if password == "" {
		err = errors.New("Password is missing")
		return
	}

	acc, err = config.Store.NewAccount(password)

	return
}
