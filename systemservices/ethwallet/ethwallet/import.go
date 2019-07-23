package ethwallet

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/engine/systemservices/ethwallet/client"
)

type importInputs struct {
	Account    encryptedKeyJSONV3 `json:"account"`
	Passphrase string             `json:"passphrase"`
}

type importOutputSuccess struct {
	Address common.Address `json:"address"`
}

func (s *Ethwallet) importA(input map[string]interface{}) (map[string]interface{}, error) {
	var inputs importInputs
	if err := client.Unmarshal(input, &inputs); err != nil {
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

	return client.Marshal(importOutputSuccess{
		Address: account.Address,
	})
}
