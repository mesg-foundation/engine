package ethwallet

import (
	"encoding/json"

	"github.com/mesg-foundation/core/client/service"
	"github.com/mesg-foundation/core/systemservices/sources/ethereum-wallet/x/xgo-ethereum/xaccounts"
)

type exportInputs struct {
	Address    string `json:"address"`
	Passphrase string `json:"passphrase"`
}

type exportOutputSuccess struct {
	Account encryptedKeyJSONV3 `json:"account"`
}

func (s *Ethwallet) export(execution *service.Execution) (string, interface{}) {
	var inputs exportInputs
	if err := execution.Data(&inputs); err != nil {
		return OutputError(err.Error())
	}

	account, err := xaccounts.GetAccount(s.keystore, inputs.Address)
	if err != nil {
		return "error", outputError{
			Message: "Account not found",
		}
	}

	keyJSON, err := s.keystore.Export(account, inputs.Passphrase, inputs.Passphrase)
	if err != nil {
		return OutputError(err.Error())
	}

	var accountJSON encryptedKeyJSONV3
	if err = json.Unmarshal(keyJSON, &accountJSON); err != nil {
		return OutputError(err.Error())
	}
	return "success", exportOutputSuccess{
		Account: accountJSON,
	}
}
