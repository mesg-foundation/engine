package ethwallet

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/core/client/service"
)

type createInputs struct {
	Passphrase string `json:"passphrase"`
}

type createOutputSuccess struct {
	Address common.Address `json:"address"`
}

func (s *Ethwallet) create(execution *service.Execution) (string, interface{}) {
	var inputs createInputs
	if err := execution.Data(&inputs); err != nil {
		return OutputError(err.Error())
	}

	account, err := s.keystore.NewAccount(inputs.Passphrase)
	if err != nil {
		return OutputError(err.Error())
	}

	return "success", createOutputSuccess{
		Address: account.Address,
	}
}
