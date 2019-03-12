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
	var output walletListOutputSuccess
	if err := p.parseResult(r, &output); err != nil {
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
	var output walletCreateOutputSuccess
	if err := p.parseResult(r, &output); err != nil {
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
	var output walletDeleteOutputSuccess
	if err := p.parseResult(r, &output); err != nil {
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
	return output, p.parseResult(r, &output)
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
	var output walletImportOutputSuccess
	if err := p.parseResult(r, &output); err != nil {
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
	var output walletImportOutputSuccess
	if err := p.parseResult(r, &output); err != nil {
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
	var output walletSignOutputSuccess
	if err := p.parseResult(r, &output); err != nil {
		return "", err
	}
	return output.SignedTransaction, nil
}

func (p *WalletProvider) parseResult(r *coreapi.ResultData, output interface{}) error {
	if r.OutputKey == "error" {
		var outputError walletErrorOutput
		if err := json.Unmarshal([]byte(r.OutputData), &outputError); err != nil {
			return err
		}
		return errors.New(outputError.Message)
	}
	return json.Unmarshal([]byte(r.OutputData), &output)
}
