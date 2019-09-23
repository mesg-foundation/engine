package ethwallet

import (
	"github.com/mesg-foundation/engine/protobuf/types"
)

func (s *Ethwallet) importA(inputs *types.Struct) (*types.Struct, error) {
	accountJSON := inputs.Fields["account"].GetStringValue()
	passphrase := inputs.Fields["passphrase"].GetStringValue()

	importedAccount, err := s.keystore.Import([]byte(accountJSON), passphrase, passphrase)
	if err != nil {
		return nil, err
	}

	return &types.Struct{Fields: map[string]*types.Value{
		"address": &types.Value{Kind: &types.Value_StringValue{importedAccount.Address.String()}},
	}}, nil
}
