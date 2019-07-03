package ethwallet

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/engine/systemservices/ethwallet/client"
)

type listOutputSuccess struct {
	Addresses []common.Address `json:"addresses"`
}

func (s *Ethwallet) list(input []byte) ([]byte, error) {
	addresses := make([]common.Address, 0)
	for _, account := range s.keystore.Accounts() {
		addresses = append(addresses, account.Address)
	}

	return client.Marshal(listOutputSuccess{
		Addresses: addresses,
	})
}
