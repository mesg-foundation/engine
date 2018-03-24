package account

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stvp/assert"
)

func TestString(t *testing.T) {
	account := &Account{
		Name:    "xxx",
		Address: common.Address{0},
		URL:     "/path/to/file",
	}
	assert.Equal(t, account.String(), "xxx 0x0000000000000000000000000000000000000000 /path/to/file")
}
