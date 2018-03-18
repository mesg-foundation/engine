package account

import (
	"testing"

	"github.com/stvp/assert"
)

func TestString(t *testing.T) {
	account := &Account{
		Name:    "xxx",
		Address: "0x00000",
	}
	assert.Equal(t, account.String(), "xxx 0x00000")
}
