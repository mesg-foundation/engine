package ethwallet

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/systemservices/ethwallet/x/xgo-ethereum/xaccounts"
)

func (s *Ethwallet) delete(inputs *types.Struct) (*types.Struct, error) {
	address := common.HexToAddress(inputs.Fields["address"].GetStringValue())
	account, err := xaccounts.GetAccount(s.keystore, address)
	if err != nil {
		return nil, errAccountNotFound
	}

	err = s.keystore.Delete(account, inputs.Fields["passphrase"].GetStringValue())
	if err != nil {
		return nil, err
	}

	return &types.Struct{
		Fields: map[string]*types.Value{
			"address": {
				Kind: &types.Value_StringValue{
					StringValue: account.Address.String(),
				},
			},
		},
	}, nil
}
