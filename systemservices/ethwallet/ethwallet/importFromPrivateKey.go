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

func (s *Ethwallet) importFromPrivateKey(execution *service.Execution) (string, interface{}) {
	var inputs importFromPrivateKeyInputs
	if err := execution.Data(&inputs); err != nil {
		return OutputError(err)
	}

	privateKey, err := crypto.HexToECDSA(inputs.PrivateKey)
	if err != nil {
		return OutputError(errors.New("cannot parse private key"))
	}

	account, err := s.keystore.ImportECDSA(privateKey, inputs.Passphrase)
	if err != nil {
		return OutputError(err)
	}

	return "success", importOutputSuccess{
		Address: account.Address,
	}
}
