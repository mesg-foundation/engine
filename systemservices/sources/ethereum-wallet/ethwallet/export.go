package ethwallet

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/core/client/service"
)

type exportInputs struct {
	Address    string `json:"address"`
	Passphrase string `json:"passphrase"`
}

type exportOutputSuccess struct {
	Account encryptedKeyJSONV3 `json:"account"`
}

func (s *Ethwallet) export(execution *service.Execution) (string, service.Data) {
	var inputs exportInputs
	if err := execution.Data(&inputs); err != nil {
		return "error", outputError{
			Message: err.Error(),
		}
	}

	_address := common.HexToAddress(inputs.Address)
	var account accounts.Account
	found := false
	for _, _account := range s.keystore.Accounts() {
		if _account.Address == _address {
			account = _account
			found = true
			break
		}
	}
	if !found {
		return "error", outputError{
			Message: "Account not found",
		}
	}

	keyJSON, err := s.keystore.Export(account, inputs.Passphrase, inputs.Passphrase)
	if err != nil {
		return "error", outputError{
			Message: err.Error(),
		}
	}

	var accountJSON encryptedKeyJSONV3
	if err = json.Unmarshal(keyJSON, &accountJSON); err != nil {
		return "error", outputError{
			Message: err.Error(),
		}
	}
	return "success", exportOutputSuccess{
		Account: accountJSON,
	}
}
