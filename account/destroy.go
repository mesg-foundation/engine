package account

import (
	"errors"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
)

// Destroy an account
func Destroy(acc accounts.Account) (err error) {
	if acc.URL.String() == "" {
		err = errors.New("Account is invalid")
	}
	err = os.Remove(acc.URL.Path)
	return
}
