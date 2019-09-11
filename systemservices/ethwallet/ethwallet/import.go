package ethwallet

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/engine/protobuf/types"
)

func (s *Ethwallet) importA(inputs *types.Struct) (*types.Struct, error) {
	account := common.HexToAddress(inputs.GetStringValue("account"))
	accountJSON, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	passphrase := inputs.GetStringValue("passphrase")
	importedAccount, err := s.keystore.Import(accountJSON, passphrase, passphrase)
	if err != nil {
		return nil, err
	}

	return types.NewStruct(map[string]*types.Value{
		"address": types.NewValueFrom(importedAccount.Address.String()),
	}), nil
}
