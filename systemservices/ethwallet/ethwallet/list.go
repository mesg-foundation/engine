package ethwallet

import (
	"github.com/mesg-foundation/engine/protobuf/types"
)

func (s *Ethwallet) list(inputs *types.Struct) (*types.Struct, error) {
	var addresses []*types.Value
	for _, account := range s.keystore.Accounts() {
		addresses = append(addresses, &types.Value{
			Kind: &types.Value_StringValue{
				StringValue: account.Address.String(),
			},
		})
	}

	return &types.Struct{
		Fields: map[string]*types.Value{
			"addresses": {
				Kind: &types.Value_ListValue{
					ListValue: &types.ListValue{
						Values: addresses,
					},
				},
			},
		},
	}, nil
}
