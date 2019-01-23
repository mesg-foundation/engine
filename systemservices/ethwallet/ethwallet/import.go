package ethwallet

import (
	"encoding/json"

	"github.com/mesg-foundation/core/client/service"
)

type importInputs struct {
	Account    encryptedKeyJSONV3 `json:"account"`
	Passphrase string             `json:"passphrase"`
}

type importOutputSuccess struct {
	Address string `json:"address"`
}

func (s *Ethwallet) importA(execution *service.Execution) (string, interface{}) {
	var inputs importInputs
	if err := execution.Data(&inputs); err != nil {
		return OutputError(err.Error())
	}

	accountJSON, err := json.Marshal(inputs.Account)
	if err != nil {
		return OutputError(err.Error())
	}

	account, err := s.keystore.Import(accountJSON, inputs.Passphrase, inputs.Passphrase)
	if err != nil {
		return OutputError(err.Error())
	}

	return "success", importOutputSuccess{
		Address: account.Address.String(),
	}
}
