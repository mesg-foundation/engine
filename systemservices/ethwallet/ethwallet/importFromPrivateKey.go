package ethwallet

import (
	"errors"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mesg-foundation/engine/protobuf/types"
)

func (s *Ethwallet) importFromPrivateKey(inputs *types.Struct) (*types.Struct, error) {
	privateKey, err := crypto.HexToECDSA(inputs.Fields["privateKey"].GetStringValue())
	if err != nil {
		return nil, errors.New("cannot parse private key")
	}

	account, err := s.keystore.ImportECDSA(privateKey, inputs.Fields["passphrase"].GetStringValue())
	if err != nil {
		return nil, err
	}

	return &types.Struct{Fields: map[string]*types.Value{
		"address": &types.Value{Kind: &types.Value_StringValue{account.Address.String()}},
	}}, nil
}
