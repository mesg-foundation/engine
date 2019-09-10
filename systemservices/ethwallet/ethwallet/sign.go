package ethwallet

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/systemservices/ethwallet/x/xgo-ethereum/xaccounts"
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

func (s *Ethwallet) sign(inputs *types.Struct) (*types.Struct, error) {
	address := common.HexToAddress(inputs.Fields["address"].GetStringValue())
	passphrase := inputs.Fields["passphrase"].GetStringValue()
	tx := inputs.Fields["transaction"].GetStructValue()

	account, err := xaccounts.GetAccount(s.keystore, address)
	if err != nil {
		return nil, errAccountNotFound
	}

	value := new(big.Int)
	if _, ok := value.SetString(tx.Fields["value"].GetStringValue(), 0); !ok {
		return nil, errCannotParseValue
	}

	gasPrice := new(big.Int)
	if _, ok := gasPrice.SetString(tx.Fields["gasPrice"].GetStringValue(), 0); !ok {
		return nil, errCannotParseGasPrice
	}

	data, err := hex.DecodeString(tx.Fields["data"].GetStringValue())
	if err != nil {
		return nil, err
	}

	transaction := ethtypes.NewTransaction(
		uint64(tx.Fields["nonce"].GetNumberValue()),
		common.HexToAddress(tx.Fields["to"].GetStringValue()),
		value,
		uint64(tx.Fields["gas"].GetNumberValue()),
		gasPrice,
		data,
	)

	signedTransaction, err := s.keystore.SignTxWithPassphrase(
		account,
		passphrase,
		transaction,
		big.NewInt(int64(tx.Fields["chainID"].GetNumberValue())),
	)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	signedTransaction.EncodeRLP(&buf)
	rawTx := fmt.Sprintf("0x%x", buf.Bytes())

	return &types.Struct{
		Fields: map[string]*types.Value{
			"signedTransaction": {
				Kind: &types.Value_StringValue{
					StringValue: rawTx,
				},
			},
		},
	}, nil
}
