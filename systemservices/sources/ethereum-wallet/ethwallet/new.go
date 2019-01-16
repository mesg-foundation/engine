package ethwallet

import (
	"github.com/mesg-foundation/core/client/service"
)

type newInputs struct {
	Passphrase string `json:"passphrase"`
}

type newOutputSuccess struct {
	Address string `json:"address"`
}

func (s *Ethwallet) new(execution *service.Execution) (string, interface{}) {
	var inputs newInputs
	if err := execution.Data(&inputs); err != nil {
		return "error", outputError{err.Error()}
	}

	account, err := s.keystore.NewAccount(inputs.Passphrase)
	if err != nil {
		return "error", outputError{err.Error()}
	}

	return "success", newOutputSuccess{
		Address: account.Address.String(),
	}
}
