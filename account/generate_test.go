package account

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts"

	"github.com/stvp/assert"
)

func TestGenerate(t *testing.T) {
	acc, err := Generate("password")

	assert.Nil(t, err)
	assert.NotEqual(t, acc, accounts.Account{})
	assert.NotEqual(t, acc.Address.String(), "")
}

func TestGeneratePasswordIsMissing(t *testing.T) {
	_, err := Generate("")
	assert.NotNil(t, err)
}
