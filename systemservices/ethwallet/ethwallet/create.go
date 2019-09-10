package ethwallet

import (
	"github.com/mesg-foundation/engine/protobuf/types"
)

func (s *Ethwallet) create(inputs *types.Struct) (*types.Struct, error) {
	account, err := s.keystore.NewAccount(inputs.Fields["passphrase"].GetStringValue())
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
