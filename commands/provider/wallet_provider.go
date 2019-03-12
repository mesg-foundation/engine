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

// List return the accounts of this wallet
func (p *WalletProvider) List() ([]string, error) {
	serviceHash, err := p.client.GetServiceHash(walletServiceKey)
	if err != nil {
		return nil, err
	}
	r, err := p.client.ExecuteAndListen(serviceHash, "list", nil)
	if err != nil {
		return nil, err
	}

	if r.OutputKey == "error" {
		var output walletErrorOutput
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

// Create creates a new account in the wallet
func (p *WalletProvider) Create(passphrase string) (string, error) {
	serviceHash, err := p.client.GetServiceHash(walletServiceKey)
	if err != nil {
		return "", err
	}
	r, err := p.client.ExecuteAndListen(serviceHash, "create", walletCreateInputs{
		Passphrase: passphrase,
	})
	if err != nil {
		return "", err
	}

	if r.OutputKey == "error" {
		var output walletErrorOutput
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

// Delete removes an account from the wallet
func (p *WalletProvider) Delete(address string, passphrase string) (string, error) {
	serviceHash, err := p.client.GetServiceHash(walletServiceKey)
	if err != nil {
		return "", err
	}
	r, err := p.client.ExecuteAndListen(serviceHash, "delete", walletDeleteInputs{
		Address:    address,
		Passphrase: passphrase,
	})
	if err != nil {
		return "", err
	}

	if r.OutputKey == "error" {
		var output walletErrorOutput
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

// Export exports an account
func (p *WalletProvider) Export(address string, passphrase string) (EncryptedKeyJSONV3, error) {
	var output EncryptedKeyJSONV3
	serviceHash, err := p.client.GetServiceHash(walletServiceKey)
	if err != nil {
		return output, err
	}
	r, err := p.client.ExecuteAndListen(serviceHash, "export", walletExportInputs{
		Address:    address,
		Passphrase: passphrase,
	})
	if err != nil {
		return output, err
	}

	if r.OutputKey == "error" {
		var outputData walletErrorOutput
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

// Import imports an account into the wallet
func (p *WalletProvider) Import(account EncryptedKeyJSONV3, passphrase string) (string, error) {
	serviceHash, err := p.client.GetServiceHash(walletServiceKey)
	if err != nil {
		return "", err
	}
	r, err := p.client.ExecuteAndListen(serviceHash, "import", &walletImportInputs{
		Account:    account,
		Passphrase: passphrase,
	})
	if err != nil {
		return "", err
	}

	if r.OutputKey == "error" {
		var output walletErrorOutput
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

// ImportFromPrivateKey imports an account from a private key
func (p *WalletProvider) ImportFromPrivateKey(privateKey string, passphrase string) (string, error) {
	serviceHash, err := p.client.GetServiceHash(walletServiceKey)
	if err != nil {
		return "", err
	}
	r, err := p.client.ExecuteAndListen(serviceHash, "importFromPrivateKey", &walletImportFromPrivateKeyInputs{
		PrivateKey: privateKey,
		Passphrase: passphrase,
	})
	if err != nil {
		return "", err
	}

	if r.OutputKey == "error" {
		var output walletErrorOutput
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

// Sign signs a transaction
func (p *WalletProvider) Sign(address string, passphrase string, transaction *Transaction) (string, error) {
	serviceHash, err := p.client.GetServiceHash(walletServiceKey)
	if err != nil {
		return "", err
	}
	r, err := p.client.ExecuteAndListen(serviceHash, "sign", &walletSignInputs{
		Address:     address,
		Passphrase:  passphrase,
		Transaction: transaction,
	})
	if err != nil {
		return "", err
	}

	if r.OutputKey == "error" {
		var output walletErrorOutput
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
