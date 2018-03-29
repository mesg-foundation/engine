package account

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mesg-foundation/application/config"
)

// Sign some data based on an account
func Sign(account accounts.Account, password string, tx *types.Transaction) (transaction *types.Transaction, err error) {
	// TODO check the last parameter
	transaction, err = config.Store.SignTxWithPassphrase(account, password, tx, big.NewInt(0))
	return
}
