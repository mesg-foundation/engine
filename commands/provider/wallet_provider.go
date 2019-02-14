package provider

import (
	"encoding/json"
	"errors"

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
		var output ErrorOutput
		if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
			return nil, err
		}
		return nil, errors.New(output.Message)
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
		var output ErrorOutput
		if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
			return common.Address{}, err
		}
		return common.Address{}, errors.New(output.Message)
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
		var output ErrorOutput
		if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
			return common.Address{}, err
		}
		return common.Address{}, errors.New(output.Message)
	}

	var output ethwalletDeleteOutputSuccess
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return common.Address{}, err
	}
	return output.Address, nil
}

func (p *WalletProvider) Export(address common.Address, passphrase string) (EncryptedKeyJSONV3, error) {
	var output EncryptedKeyJSONV3
	r, err := p.client.ExecuteAndListen(walletServiceID, "export", ethwalletExportInputs{
		Address:    address,
		Passphrase: passphrase,
	})
	if err != nil {
		return output, err
	}

	if r.OutputKey == "error" {
		var outputData ErrorOutput
		if err := json.Unmarshal([]byte(r.OutputData), &outputData); err != nil {
			return output, err
		}
		return output, errors.New(outputData.Message)
	}

	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return output, err
	}
	return output, nil
}

func (p *WalletProvider) Import(account EncryptedKeyJSONV3, passphrase string) (common.Address, error) {
	r, err := p.client.ExecuteAndListen(walletServiceID, "import", &ethwalletImportInputs{
		Account:    account,
		Passphrase: passphrase,
	})
	if err != nil {
		return common.Address{}, err
	}

	if r.OutputKey == "error" {
		var output ErrorOutput
		if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
			return common.Address{}, err
		}
		return common.Address{}, errors.New(output.Message)
	}

	var output ethwalletImportOutputSuccess
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return common.Address{}, err
	}
	return output.Address, nil
}

func (p *WalletProvider) ImportFromPrivateKey(privateKey string, passphrase string) (common.Address, error) {
	r, err := p.client.ExecuteAndListen(walletServiceID, "importFromPrivateKey", &ethwalletImportFromPrivateKeyInputs{
		PrivateKey: privateKey,
		Passphrase: passphrase,
	})
	if err != nil {
		return common.Address{}, err
	}

	if r.OutputKey == "error" {
		var output ErrorOutput
		if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
			return common.Address{}, err
		}
		return common.Address{}, errors.New(output.Message)
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
		var output ErrorOutput
		if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
			return "", err
		}
		return "", errors.New(output.Message)
	}

	var output ethwalletSignOutputSuccess
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return "", err
	}
	return output.SignedTransaction, nil
}
