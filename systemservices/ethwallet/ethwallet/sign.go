package ethwallet

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mesg-foundation/core/client/service"
	"github.com/mesg-foundation/core/systemservices/ethwallet/x/xgo-ethereum/xaccounts"
)

type signInputs struct {
	Address     common.Address `json:"address"`
	Passphrase  string         `json:"passphrase"`
	Transaction *transaction   `json:"transaction"`
}

type transaction struct {
	ChainID  int64          `json:"chainID"`
	Nonce    uint64         `json:"nonce"`
	To       common.Address `json:"to"`
	Value    string         `json:"value"`
	Gas      uint64         `json:"gas"`
	GasPrice string         `json:"gasPrice"`
	Data     hexutil.Bytes  `json:"data"`
}

type signOutputSuccess struct {
	SignedTransaction string `json:"signedTransaction"`
}

func (s *Ethwallet) sign(execution *service.Execution) (interface{}, error) {
	var inputs signInputs
	if err := execution.Data(&inputs); err != nil {
		return nil, err
	}

	account, err := xaccounts.GetAccount(s.keystore, inputs.Address)
	if err != nil {
		return nil, errAccountNotFound
	}

	value := new(big.Int)
	if _, ok := value.SetString(inputs.Transaction.Value, 0); !ok {
		return nil, errCannotParseValue
	}

	gasPrice := new(big.Int)
	if _, ok := gasPrice.SetString(inputs.Transaction.GasPrice, 0); !ok {
		return nil, errCannotParseGasPrice
	}

	transaction := types.NewTransaction(inputs.Transaction.Nonce, inputs.Transaction.To, value, inputs.Transaction.Gas, gasPrice, inputs.Transaction.Data)

	signedTransaction, err := s.keystore.SignTxWithPassphrase(account, inputs.Passphrase, transaction, big.NewInt(inputs.Transaction.ChainID))
	if err != nil {
		return nil, err
	}

	var buff bytes.Buffer
	signedTransaction.EncodeRLP(&buff)
	rawTx := fmt.Sprintf("0x%x", buff.Bytes())

	return signOutputSuccess{
		SignedTransaction: rawTx,
	}, nil
}
