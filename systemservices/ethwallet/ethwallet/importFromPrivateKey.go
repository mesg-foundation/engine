package ethwallet

import (
	"errors"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mesg-foundation/engine/protobuf/types"
)

func (s *Ethwallet) importFromPrivateKey(inputs *types.Struct) (*types.Struct, error) {
	privateKey, err := crypto.HexToECDSA(inputs.GetStringValue("privateKey"))
	if err != nil {
		return nil, errors.New("cannot parse private key")
	}

	account, err := s.keystore.ImportECDSA(privateKey, inputs.GetStringValue("passphrase"))
	if err != nil {
		return nil, err
	}

	return types.NewStruct(map[string]*types.Value{
		"address": types.NewValueFrom(account.Address.String()),
	}), nil
}
