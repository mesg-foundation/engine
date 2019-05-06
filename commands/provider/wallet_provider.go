package provider

import (
	"encoding/json"

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
	var output walletListOutputSuccess
	return output.Addresses, p.call("list", nil, &output)
}

// Create creates a new account in the wallet
func (p *WalletProvider) Create(passphrase string) (string, error) {
	input := walletCreateInputs{
		Passphrase: passphrase,
	}
	var output walletCreateOutputSuccess
	return output.Address, p.call("create", &input, &output)
}

// Delete removes an account from the wallet
func (p *WalletProvider) Delete(address string, passphrase string) (string, error) {
	input := walletDeleteInputs{
		Address:    address,
		Passphrase: passphrase,
	}
	var output walletDeleteOutputSuccess
	return output.Address, p.call("delete", &input, &output)
}

// Export exports an account
func (p *WalletProvider) Export(address string, passphrase string) (WalletEncryptedKeyJSONV3, error) {
	input := walletExportInputs{
		Address:    address,
		Passphrase: passphrase,
	}
	var output WalletEncryptedKeyJSONV3
	return output, p.call("export", &input, &output)
}

// Import imports an account into the wallet
func (p *WalletProvider) Import(account WalletEncryptedKeyJSONV3, passphrase string) (string, error) {
	input := walletImportInputs{
		Account:    account,
		Passphrase: passphrase,
	}
	var output walletImportOutputSuccess
	return output.Address, p.call("import", &input, &output)
}

// ImportFromPrivateKey imports an account from a private key
func (p *WalletProvider) ImportFromPrivateKey(privateKey string, passphrase string) (string, error) {
	input := walletImportFromPrivateKeyInputs{
		PrivateKey: privateKey,
		Passphrase: passphrase,
	}
	var output walletImportOutputSuccess
	return output.Address, p.call("importFromPrivateKey", &input, &output)

}

// Sign signs a transaction
func (p *WalletProvider) Sign(address string, passphrase string, transaction Transaction) (string, error) {
	input := walletSignInputs{
		Address:     address,
		Passphrase:  passphrase,
		Transaction: transaction,
	}
	var output walletSignOutputSuccess
	return output.SignedTransaction, p.call("sign", &input, &output)
}

func (p *WalletProvider) call(task string, inputs interface{}, output interface{}) error {
	serviceHash, err := p.client.GetServiceHash(walletServiceKey)
	if err != nil {
		return err
	}
	r, err := p.client.ExecuteAndListen(serviceHash, task, &inputs)
	if err != nil {
		return err
	}
	return p.parseResult(r, &output)
}

func (p *WalletProvider) parseResult(r *coreapi.ResultData, output interface{}) error {
	if r.OutputKey == "error" {
		var outputError walletErrorOutput
		if err := json.Unmarshal([]byte(r.OutputData), &outputError); err != nil {
			return err
		}
		return outputError
	}
	return json.Unmarshal([]byte(r.OutputData), &output)
}
