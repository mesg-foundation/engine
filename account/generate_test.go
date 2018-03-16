package account

import (
	"testing"

	"github.com/stvp/assert"
)

func TestGenerate(t *testing.T) {
	addr, seed, err := generate("password", "name")

	// err
	assert.Nil(t, err, "The error should be nil")

	// addr
	assert.NotNil(t, addr, "The address should not be nil")
	assert.True(t, addr != "", "The address should not be empty")

	// seed
	assert.NotNil(t, seed, "The seed should not be nil")
	assert.True(t, seed != "", "The seed should not be empty")
}

func TestGenerateNameIsMissing(t *testing.T) {
	_, _, err := generate("password", "")
	assert.NotNil(t, err, "The error should not be nil")
	assert.Equal(t, err.Error(), "Name is missing")
}

func TestGeneratePasswordIsMissing(t *testing.T) {
	_, _, err := generate("", "name")
	assert.NotNil(t, err, "The error should not be nil")
	assert.Equal(t, err.Error(), "Password is missing")
}

func TestGenerateFromAccount(t *testing.T) {
	account := Account{
		Name:     "Testx",
		Password: "xxxxx",
	}
	account.Generate()
	assert.NotEqual(t, account.Address, "")
	assert.NotEqual(t, account.Seed, "")
}

func TestGenerateFromAccountWithNoPassword(t *testing.T) {
	account := Account{
		Name: "Testx",
	}
	err := account.Generate()
	assert.NotNil(t, err)
}
