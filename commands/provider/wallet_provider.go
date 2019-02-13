package provider

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mesg-foundation/core/protobuf/coreapi"
)

// WalletProvider is a struct that provides all methods required by wallet command.
type WalletProvider struct {
	client client
}

// NewWalletProvider creates new WalletProvider.
func NewWalletProvider(c coreapi.CoreClient) *WalletProvider {
	return &WalletProvider{client: client{c}}
}

func (p *WalletProvider) List() ([]common.Address, error) {
	r, err := p.client.ExecuteAndListen(walletServiceID, "list", nil)
	if err != nil {
		return nil, err
	}
	if r.OutputKey == "error" {
		return nil, fmt.Errorf("servcie %s task %s return error: %s", walletServiceID, r.TaskKey, r.OutputData)
	}

	var output ethwalletListOutputSuccess
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return nil, err
	}
	return output.Addresses, nil
}

func (p *WalletProvider) Create(passphrase string) (common.Address, error) {
	r, err := p.client.ExecuteAndListen(walletServiceID, "create", ethwalletCreateInputs{
		Passphrase: passphrase,
	})
	if err != nil {
		return common.Address{}, err
	}
	if r.OutputKey == "error" {
		return common.Address{}, fmt.Errorf("servcie %s task %s return error: %s", walletServiceID, r.TaskKey, r.OutputData)
	}

	var output ethwalletCreateOutputSuccess
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return common.Address{}, err
	}
	return output.Address, nil
}

func (p *WalletProvider) Delete(address common.Address, passphrase string) (common.Address, error) {
	r, err := p.client.ExecuteAndListen(walletServiceID, "delete", ethwalletDeleteInputs{
		Address:    address,
		Passphrase: passphrase,
	})
	if err != nil {
		return common.Address{}, err
	}
	if r.OutputKey == "error" {
		return common.Address{}, fmt.Errorf("servcie %s task %s return error: %s", walletServiceID, r.TaskKey, r.OutputData)
	}

	var output ethwalletDeleteOutputSuccess
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return common.Address{}, err
	}
	return output.Address, nil
}

func (p *WalletProvider) Export(address common.Address, passphrase string) (EncryptedKeyJSONV3, error) {
	r, err := p.client.ExecuteAndListen(walletServiceID, "delete", ethwalletExportInputs{
		Address:    address,
		Passphrase: passphrase,
	})
	if err != nil {
		return EncryptedKeyJSONV3{}, err
	}
	if r.OutputKey == "error" {
		return EncryptedKeyJSONV3{}, fmt.Errorf("servcie %s task %s return error: %s", walletServiceID, r.TaskKey, r.OutputData)
	}

	var output EncryptedKeyJSONV3
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return EncryptedKeyJSONV3{}, err
	}
	return output, nil
}
func (p *WalletProvider) Import(address common.Address, passphrase string, account EncryptedKeyJSONV3) (common.Address, error) {
	r, err := p.client.ExecuteAndListen(walletServiceID, "sign", &ethwalletImportInputs{
		Account:    account,
		Passphrase: passphrase,
	})
	if err != nil {
		return common.Address{}, err
	}
	if r.OutputKey == "error" {
		return common.Address{}, fmt.Errorf("servcie %s task %s return error: %s", walletServiceID, r.TaskKey, r.OutputData)
	}

	var output ethwalletImportOutputSuccess
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return common.Address{}, err
	}
	return output.Address, nil
}
func (p *WalletProvider) Sign(address common.Address, passphrase string, transaction *Transaction) (string, error) {
	r, err := p.client.ExecuteAndListen(walletServiceID, "sign", &ethwalletSignInputs{
		Address:     address,
		Passphrase:  passphrase,
		Transaction: transaction,
	})
	if err != nil {
		return "", err
	}
	if r.OutputKey == "error" {
		return "", fmt.Errorf("servcie %s task %s return error: %s", walletServiceID, r.TaskKey, r.OutputData)
	}

	var output ethwalletSignOutputSuccess
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return "", err
	}
	return output.SignedTransaction, nil
}
