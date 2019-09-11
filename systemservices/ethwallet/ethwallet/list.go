package ethwallet

import (
	"github.com/mesg-foundation/engine/protobuf/types"
)

func (s *Ethwallet) list(inputs *types.Struct) (*types.Struct, error) {
	var addresses []*types.Value
	for _, account := range s.keystore.Accounts() {
		addresses = append(addresses, types.NewValueFrom(account.Address.String()))
	}

	return types.NewStruct(map[string]*types.Value{
		"addresses": types.NewValueFrom(addresses),
	}), nil
}
