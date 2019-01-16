package ethwallet

import (
	"github.com/mesg-foundation/core/client/service"
)

type listOutputSuccess struct {
	Addresses []string `json:"addresses"`
}

func (s *Ethwallet) list(execution *service.Execution) (string, service.Data) {
	addresses := make([]string, 0)
	for _, account := range s.keystore.Accounts() {
		addresses = append(addresses, account.Address.String())
	}

	return "success", listOutputSuccess{
		Addresses: addresses,
	}
}
