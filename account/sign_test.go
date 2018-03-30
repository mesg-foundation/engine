package account

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stvp/assert"
)

func TestSign(t *testing.T) {
	account, _ := Generate("TestSign")
	amount := big.NewInt(0)
	gasPrice := big.NewInt(0)
	tx := types.NewTransaction(0, account.Address, amount, 21000, gasPrice, []byte(""))
	res, err := signTx(account, "TestSign", tx)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	Destroy(account)
}
