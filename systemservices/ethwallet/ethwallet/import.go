package ethwallet

import (
	"github.com/mesg-foundation/engine/protobuf/types"
)

func (s *Ethwallet) importA(inputs *types.Struct) (*types.Struct, error) {
	accountJSON := inputs.GetStringValue("account")
	passphrase := inputs.GetStringValue("passphrase")

	importedAccount, err := s.keystore.Import([]byte(accountJSON), passphrase, passphrase)
	if err != nil {
		return nil, err
	}

	return types.NewStruct(map[string]*types.Value{
		"address": types.NewValueFrom(importedAccount.Address.String()),
	}), nil
}
