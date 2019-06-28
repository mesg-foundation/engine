package ethwallet

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/engine/client/service"
)

type listOutputSuccess struct {
	Addresses []common.Address `json:"addresses"`
}

func (s *Ethwallet) list(execution *service.Execution) (interface{}, error) {
	addresses := make([]common.Address, 0)
	for _, account := range s.keystore.Accounts() {
		addresses = append(addresses, account.Address)
	}

	return listOutputSuccess{
		Addresses: addresses,
	}, nil
}
