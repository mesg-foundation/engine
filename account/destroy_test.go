package account

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts"

	"github.com/stvp/assert"
)

func TestDestroy(t *testing.T) {
	acc, _ := Generate("test")
	err := Destroy(acc)
	assert.Nil(t, err)
}

func TestDestroyInvalidAccount(t *testing.T) {
	acc := accounts.Account{}
	err := Destroy(acc)
	assert.NotNil(t, err)
}
