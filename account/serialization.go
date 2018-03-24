package account

import (
	"github.com/ethereum/go-ethereum/accounts"
)

// Import a file and return a Account object
func Import(filePath string, name string) (account accounts.Account, err error) {
	if name == "" {
		name = "Test A"
	}
	// TODO add import
	// account = &Account{
	// 	Address: common.Address{0},
	// 	Name:    name,
	// }
	return
}
