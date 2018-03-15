package account

import (
	"errors"
)

// TODO add real account creation
func generate(password string, name string) (string, string, error) {
	var addr, seed string
	if password == "" {
		return addr, seed, errors.New("Password is missing")
	}
	if name == "" {
		return addr, seed, errors.New("Name is missing")
	}
	addr = "0x0000000000000000000000000000000000000000"
	seed = "this is my long secure seed that help me regenerate my account keys"
	return addr, seed, nil
}

// Generate an account based on some predefined data
func (account *Account) Generate() error {
	addr, seed, err := generate(account.Password, account.Name)
	if err != nil {
		return err
	}
	account.Address = addr
	account.Seed = seed
	return nil
}
