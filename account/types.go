package account

import "github.com/ethereum/go-ethereum/common"

// Account is a structure that contains all information about an account
type Account struct {
	Name     string
	Address  common.Address
	Password string
	Seed     string
}

func (account *Account) String() (desc string) {
	return account.Name + " " + account.Address.String()
}
