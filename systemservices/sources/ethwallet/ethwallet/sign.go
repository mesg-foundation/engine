package ethwallet

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mesg-foundation/core/client/service"
	"github.com/mesg-foundation/core/systemservices/sources/ethwallet/x/xgo-ethereum/xaccounts"
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

func (s *Ethwallet) sign(execution *service.Execution) (string, interface{}) {
	var inputs signInputs
	if err := execution.Data(&inputs); err != nil {
		return OutputError(err.Error())
	}

	account, err := xaccounts.GetAccount(s.keystore, inputs.Address)
	if err != nil {
		return "error", outputError{
			Message: "Account not found",
		}
	}

	transaction := types.NewTransaction(inputs.Transaction.Nonce, inputs.Transaction.To, inputs.Transaction.Value, inputs.Transaction.GasLimit, inputs.Transaction.GasPrice, inputs.Transaction.Data)

	signedTransaction, err := s.keystore.SignTxWithPassphrase(account, inputs.Passphrase, transaction, big.NewInt(inputs.ChainID))
	if err != nil {
		return OutputError(err.Error())
	}

	return "success", signOutputSuccess{
		SignedTransaction: signedTransaction,
	}
}
