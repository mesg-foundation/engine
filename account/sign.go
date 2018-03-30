package account

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mesg-foundation/application/config"
)

// sign some data based on an account
func signTx(account accounts.Account, password string, tx *types.Transaction) (transaction *types.Transaction, err error) {
	// TODO check the last parameter
	transaction, err = config.Store.SignTxWithPassphrase(account, password, tx, config.NetworkID)
	return
}
