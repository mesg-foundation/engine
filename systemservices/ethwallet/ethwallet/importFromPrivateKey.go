package ethwallet

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mesg-foundation/core/client/service"
)

type importFromPrivateKeyInputs struct {
	PrivateKey string `json:"privateKey"`
	Passphrase string `json:"passphrase"`
}

type importFromPrivateKeyOutputSuccess struct {
	Address common.Address `json:"address"`
}

func (s *Ethwallet) importFromPrivateKey(execution *service.Execution) (interface{}, error) {
	var inputs importFromPrivateKeyInputs
	if err := execution.Data(&inputs); err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(inputs.PrivateKey)
	if err != nil {
		return nil, errors.New("cannot parse private key")
	}

	account, err := s.keystore.ImportECDSA(privateKey, inputs.Passphrase)
	if err != nil {
		return nil, err
	}

	return importOutputSuccess{
		Address: account.Address,
	}, nil
}
