package ethwallet

import (
	"github.com/mesg-foundation/core/client/service"
	"github.com/mesg-foundation/core/systemservices/sources/ethereum-wallet/x/xgo-ethereum/xaccounts"
)

type deleteInputs struct {
	Address    string `json:"address"`
	Passphrase string `json:"passphrase"`
}

func (s *Ethwallet) delete(execution *service.Execution) (string, interface{}) {
	var inputs deleteInputs
	if err := execution.Data(&inputs); err != nil {
		return OutputError(err.Error())
	}

	account, err := xaccounts.GetAccount(s.keystore, inputs.Address)
	if err != nil {
		return "error", outputError{
			Message: "Account not found",
		}
	}

	if err := s.keystore.Delete(account, inputs.Passphrase); err != nil {
		return OutputError(err.Error())
	}

	return "success", nil
}
