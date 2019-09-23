package ethwallet

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/systemservices/ethwallet/x/xgo-ethereum/xaccounts"
)

func (s *Ethwallet) sign(inputs *types.Struct) (*types.Struct, error) {
	address := common.HexToAddress(inputs.Fields["address"].GetStringValue())
	passphrase := inputs.Fields["passphrase"].GetStringValue()
	tx := inputs.Fields["transaction"].GetStructValue()

	account, err := xaccounts.GetAccount(s.keystore, address)
	if err != nil {
		return nil, err
	}

	value := new(big.Int)
	if _, ok := value.SetString(tx.Fields["value"].GetStringValue(), 0); !ok {
		return nil, errors.New("cannot parse value")
	}

	gasPrice := new(big.Int)
	if _, ok := gasPrice.SetString(tx.Fields["gasPrice"].GetStringValue(), 0); !ok {
		return nil, errors.New("cannot parse gasPrice")
	}

	data, err := hexutil.Decode(tx.Fields["data"].GetStringValue())
	if err != nil {
		return nil, fmt.Errorf("cannot parse data: %w", err)
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

	return &types.Struct{Fields: map[string]*types.Value{
		"signedTransaction": &types.Value{Kind: &types.Value_StringValue{rawTx}},
	}}, nil
}
