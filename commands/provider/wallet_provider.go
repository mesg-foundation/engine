package provider

import (
	"encoding/json"
	"errors"

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

func (p *WalletProvider) List() ([]string, error) {
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

	var output walletListOutputSuccess
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return nil, err
	}
	return output.Addresses, nil
}

func (p *WalletProvider) Create(passphrase string) (string, error) {
	r, err := p.client.ExecuteAndListen(walletServiceID, "create", walletCreateInputs{
		Passphrase: passphrase,
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

	var output walletCreateOutputSuccess
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return "", err
	}
	return output.Address, nil
}

func (p *WalletProvider) Delete(address string, passphrase string) (string, error) {
	r, err := p.client.ExecuteAndListen(walletServiceID, "delete", walletDeleteInputs{
		Address:    address,
		Passphrase: passphrase,
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

	var output walletDeleteOutputSuccess
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return "", err
	}
	return output.Address, nil
}

func (p *WalletProvider) Export(address string, passphrase string) (EncryptedKeyJSONV3, error) {
	var output EncryptedKeyJSONV3
	r, err := p.client.ExecuteAndListen(walletServiceID, "export", walletExportInputs{
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

func (p *WalletProvider) Import(account EncryptedKeyJSONV3, passphrase string) (string, error) {
	r, err := p.client.ExecuteAndListen(walletServiceID, "import", &walletImportInputs{
		Account:    account,
		Passphrase: passphrase,
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

	var output walletImportOutputSuccess
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return "", err
	}
	return output.Address, nil
}

func (p *WalletProvider) ImportFromPrivateKey(privateKey string, passphrase string) (string, error) {
	r, err := p.client.ExecuteAndListen(walletServiceID, "importFromPrivateKey", &walletImportFromPrivateKeyInputs{
		PrivateKey: privateKey,
		Passphrase: passphrase,
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

	var output walletImportOutputSuccess
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return "", err
	}
	return output.Address, nil
}

func (p *WalletProvider) Sign(address string, passphrase string, transaction *Transaction) (string, error) {
	r, err := p.client.ExecuteAndListen(walletServiceID, "sign", &walletSignInputs{
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

	var output walletSignOutputSuccess
	if err := json.Unmarshal([]byte(r.OutputData), &output); err != nil {
		return "", err
	}
	return output.SignedTransaction, nil
}
