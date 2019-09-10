package ethwallet

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/engine/protobuf/types"
)

func (s *Ethwallet) importA(inputs *types.Struct) (*types.Struct, error) {
	account := common.HexToAddress(inputs.Fields["account"].GetStringValue())
	accountJSON, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	passphrase := inputs.Fields["passphrase"].GetStringValue()
	importedAccount, err := s.keystore.Import(accountJSON, passphrase, passphrase)
	if err != nil {
		return nil, err
	}

	return &types.Struct{
		Fields: map[string]*types.Value{
			"address": {
				Kind: &types.Value_StringValue{
					StringValue: importedAccount.Address.String(),
				},
			},
		},
	}, nil
}
