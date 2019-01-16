package ethwallet

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/core/client/service"
)

type deleteInputs struct {
	Address    string `json:"address"`
	Passphrase string `json:"passphrase"`
}

func (s *Ethwallet) delete(execution *service.Execution) (string, interface{}) {
	var inputs deleteInputs
	if err := execution.Data(&inputs); err != nil {
		return "error", outputError{
			Message: err.Error(),
		}
	}

	_address := common.HexToAddress(inputs.Address)
	var account accounts.Account
	found := false
	for _, _account := range s.keystore.Accounts() {
		if _account.Address == _address {
			account = _account
			found = true
			break
		}
	}
	if !found {
		return "error", outputError{
			Message: "Account not found",
		}
	}

	if err := s.keystore.Delete(account, inputs.Passphrase); err != nil {
		return "error", outputError{
			Message: err.Error(),
		}
	}

	return "success", nil
}
