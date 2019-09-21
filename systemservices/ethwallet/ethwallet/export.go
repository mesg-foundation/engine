package ethwallet

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/systemservices/ethwallet/x/xgo-ethereum/xaccounts"
)

func (s *Ethwallet) export(inputs *types.Struct) (*types.Struct, error) {
	address := common.HexToAddress(inputs.GetStringValue("address"))
	passphrase := inputs.GetStringValue("passphrase")

	account, err := xaccounts.GetAccount(s.keystore, address)
	if err != nil {
		return nil, errAccountNotFound
	}

	accountJSON, err := s.keystore.Export(account, passphrase, passphrase)
	if err != nil {
		return nil, err
	}

	return types.NewStruct(map[string]*types.Value{
		"account": types.NewValueFrom(string(accountJSON)),
	}), nil
}
