package ethwallet

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/core/client/service"
	"github.com/mesg-foundation/core/systemservices/ethwallet/x/xgo-ethereum/xaccounts"
)

type deleteInputs struct {
	Address    common.Address `json:"address"`
	Passphrase string         `json:"passphrase"`
}

type deleteOutputSuccess struct {
	Address common.Address `json:"address"`
}

func (s *Ethwallet) delete(execution *service.Execution) (string, interface{}) {
	var inputs deleteInputs
	if err := execution.Data(&inputs); err != nil {
		return OutputError(err.Error())
	}

	account, err := xaccounts.GetAccount(s.keystore, inputs.Address)
	if err != nil {
		return OutputError("Account not found")
	}

	if err := s.keystore.Delete(account, inputs.Passphrase); err != nil {
		return OutputError(err.Error())
	}

	return "success", deleteOutputSuccess{
		Address: account.Address,
	}
}
