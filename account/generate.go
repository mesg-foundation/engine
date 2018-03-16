package account

import (
	"errors"
)

// TODO add real account creation
func generate(password string, name string) (addr string, seed string, err error) {
	if password == "" {
		err = errors.New("Password is missing")
	}
	if name == "" {
		err = errors.New("Name is missing")
	}
	if err != nil {
		return
	}
	addr = "0x0000000000000000000000000000000000000000"
	seed = "this is my long secure seed that help me regenerate my account keys"
	return
}

// Generate an account based on some predefined data
func (account *Account) Generate() (err error) {
	addr, seed, err := generate(account.Password, account.Name)
	if err != nil {
		return
	}
	account.Address = addr
	account.Seed = seed
	return
}
