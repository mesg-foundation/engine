package ethwallet

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/engine/systemservices/ethwallet/client"
)

type createInputs struct {
	Passphrase string `json:"passphrase"`
}

type createOutputSuccess struct {
	Address common.Address `json:"address"`
}

func (s *Ethwallet) create(input map[string]interface{}) (map[string]interface{}, error) {
	var inputs createInputs
	if err := client.Unmarshal(input, &inputs); err != nil {
		return nil, err
	}

	account, err := s.keystore.NewAccount(inputs.Passphrase)
	if err != nil {
		return nil, err
	}

	return client.Marshal(createOutputSuccess{
		Address: account.Address,
	})
}
