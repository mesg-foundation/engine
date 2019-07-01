package ethwallet

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/engine/systemservices/ethwallet/client"
	"github.com/mesg-foundation/engine/systemservices/ethwallet/x/xgo-ethereum/xaccounts"
)

type deleteInputs struct {
	Address    common.Address `json:"address"`
	Passphrase string         `json:"passphrase"`
}

type deleteOutputSuccess struct {
	Address common.Address `json:"address"`
}

func (s *Ethwallet) delete(input []byte) ([]byte, error) {
	var inputs deleteInputs
	if err := client.Unmarshal(input, &inputs); err != nil {
		return nil, err
	}

	account, err := xaccounts.GetAccount(s.keystore, inputs.Address)
	if err != nil {
		return nil, errAccountNotFound
	}

	if err := s.keystore.Delete(account, inputs.Passphrase); err != nil {
		return nil, err
	}

	return client.Marshal(deleteOutputSuccess{
		Address: account.Address,
	})
}
