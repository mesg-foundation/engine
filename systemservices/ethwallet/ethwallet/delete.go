package ethwallet

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/systemservices/ethwallet/x/xgo-ethereum/xaccounts"
)

func (s *Ethwallet) delete(inputs *types.Struct) (*types.Struct, error) {
	address := common.HexToAddress(inputs.GetStringValue("address"))
	account, err := xaccounts.GetAccount(s.keystore, address)
	if err != nil {
		return nil, errAccountNotFound
	}

	err = s.keystore.Delete(account, inputs.GetStringValue("passphrase"))
	if err != nil {
		return nil, err
	}

	return types.NewStruct(map[string]*types.Value{
		"address": types.NewValueFrom(account.Address.String()),
	}), nil
}
