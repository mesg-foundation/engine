package ethwallet

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mesg-foundation/core/client/service"
)

type signInputs struct {
	Address     string       `json:"address"`
	Passphrase  string       `json:"passphrase"`
	Transaction *transaction `json:"transaction"`
	ChainID     int64        `json:"chainID"`
}

type transaction struct {
	Nonce    uint64         `json:"nonce"`
	To       common.Address `json:"to"`
	Value    *big.Int       `json:"value"`
	GasLimit uint64         `json:"gasLimit"`
	GasPrice *big.Int       `json:"gasPrice"`
	Data     []byte         `json:"data"`
}

type signOutputSuccess struct {
	SignedTransaction *types.Transaction `json:"signedTransaction"`
}

func (s *Ethwallet) sign(execution *service.Execution) (string, service.Data) {
	var inputs signInputs
	if err := execution.Data(&inputs); err != nil {
		return "error", outputError{
			Message: err.Error(),
		}
	}

	// TODO: refacto following block with export
	_address := common.HexToAddress(inputs.Address)
	var account accounts.Account
	found := false
	for _, _account := range s.keystore.Accounts() {
		if _account.Address == _address {
			account = _account
			found = true
			break
		}
	}
	if !found {
		return "error", outputError{
			Message: "Account not found",
		}
	}

	transaction := types.NewTransaction(inputs.Transaction.Nonce, inputs.Transaction.To, inputs.Transaction.Value, inputs.Transaction.GasLimit, inputs.Transaction.GasPrice, inputs.Transaction.Data)

	signedTransaction, err := s.keystore.SignTxWithPassphrase(account, inputs.Passphrase, transaction, big.NewInt(inputs.ChainID))
	if err != nil {
		return "error", outputError{
			Message: err.Error(),
		}
	}

	return "success", signOutputSuccess{
		SignedTransaction: signedTransaction,
	}
}
