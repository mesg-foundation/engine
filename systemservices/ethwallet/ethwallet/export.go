package ethwallet

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/core/client/service"
	"github.com/mesg-foundation/core/systemservices/ethwallet/x/xgo-ethereum/xaccounts"
)

type exportInputs struct {
	Address    common.Address `json:"address"`
	Passphrase string         `json:"passphrase"`
}

func (s *Ethwallet) export(execution *service.Execution) (string, interface{}) {
	var inputs exportInputs
	if err := execution.Data(&inputs); err != nil {
		return OutputError(err)
	}

	account, err := xaccounts.GetAccount(s.keystore, inputs.Address)
	if err != nil {
		return OutputError(errAccountNotFound)
	}

	keyJSON, err := s.keystore.Export(account, inputs.Passphrase, inputs.Passphrase)
	if err != nil {
		return OutputError(err)
	}

	var accountJSON encryptedKeyJSONV3
	if err = json.Unmarshal(keyJSON, &accountJSON); err != nil {
		return OutputError(err)
	}
	return "success", accountJSON
}
