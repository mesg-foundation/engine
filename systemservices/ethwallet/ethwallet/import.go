package ethwallet

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/core/client/service"
)

type importInputs struct {
	Account    encryptedKeyJSONV3 `json:"account"`
	Passphrase string             `json:"passphrase"`
}

type importOutputSuccess struct {
	Address common.Address `json:"address"`
}

func (s *Ethwallet) importA(execution *service.Execution) (interface{}, error) {
	var inputs importInputs
	if err := execution.Data(&inputs); err != nil {
		return nil, err
	}

	accountJSON, err := json.Marshal(inputs.Account)
	if err != nil {
		return nil, err
	}

	account, err := s.keystore.Import(accountJSON, inputs.Passphrase, inputs.Passphrase)
	if err != nil {
		return nil, err
	}

	return importOutputSuccess{
		Address: account.Address,
	}, nil
}
