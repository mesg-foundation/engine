package account

import (
	"testing"

	"github.com/stvp/assert"
)

func TestGenerate(t *testing.T) {
	addr, err := generate("password", "name")

	assert.Nil(t, err)
	assert.NotNil(t, addr)
	assert.NotEqual(t, addr.String(), "")
	assert.Equal(t, len(addr.String()), 42)
}

func TestGenerateNameIsMissing(t *testing.T) {
	_, err := generate("password", "")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Name is missing")
}

func TestGeneratePasswordIsMissing(t *testing.T) {
	_, err := generate("", "name")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Password is missing")
}

func TestGenerateFromAccount(t *testing.T) {
	account := Account{
		Name:     "Testx",
		Password: "xxxxx",
	}
	account.Generate()
	assert.NotEqual(t, account.Address, "")
}

func TestGenerateFromAccountWithNoPassword(t *testing.T) {
	account := Account{
		Name: "Testx",
	}
	err := account.Generate()
	assert.NotNil(t, err)
}
