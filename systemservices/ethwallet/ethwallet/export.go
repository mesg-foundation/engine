package ethwallet

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/engine/systemservices/ethwallet/client"
	"github.com/mesg-foundation/engine/systemservices/ethwallet/x/xgo-ethereum/xaccounts"
)

type exportInputs struct {
	Address    common.Address `json:"address"`
	Passphrase string         `json:"passphrase"`
}

func (s *Ethwallet) export(input map[string]interface{}) (map[string]interface{}, error) {
	var inputs exportInputs
	if err := client.Unmarshal(input, &inputs); err != nil {
		return nil, err
	}

	account, err := xaccounts.GetAccount(s.keystore, inputs.Address)
	if err != nil {
		return nil, errAccountNotFound
	}

	keyJSON, err := s.keystore.Export(account, inputs.Passphrase, inputs.Passphrase)
	if err != nil {
		return nil, err
	}

	var accountJSON encryptedKeyJSONV3
	if err = json.Unmarshal(keyJSON, &accountJSON); err != nil {
		return nil, err
	}
	return client.Marshal(accountJSON)
}
