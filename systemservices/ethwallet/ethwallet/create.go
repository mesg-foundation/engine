package ethwallet

import (
	"github.com/mesg-foundation/engine/protobuf/types"
)

func (s *Ethwallet) create(inputs *types.Struct) (*types.Struct, error) {
	account, err := s.keystore.NewAccount(inputs.GetStringValue("passphrase"))
	if err != nil {
		return nil, err
	}

	return types.NewStruct(map[string]*types.Value{
		"address": types.NewValueFrom(account.Address.String()),
	}), nil
}
