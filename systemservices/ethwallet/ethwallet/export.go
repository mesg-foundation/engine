package ethwallet

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/systemservices/ethwallet/x/xgo-ethereum/xaccounts"
)

func (s *Ethwallet) export(inputs *types.Struct) (*types.Struct, error) {
	address := common.HexToAddress(inputs.Fields["address"].GetStringValue())
	passphrase := inputs.Fields["passphrase"].GetStringValue()

	account, err := xaccounts.GetAccount(s.keystore, address)
	if err != nil {
		return nil, err
	}

	accountJSON, err := s.keystore.Export(account, passphrase, passphrase)
	if err != nil {
		return nil, err
	}

	return &types.Struct{Fields: map[string]*types.Value{
		"account": &types.Value{Kind: &types.Value_StringValue{string(accountJSON)}},
	}}, nil
}
