package provider

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mesg-foundation/core/protobuf/coreapi"
)

// WalletProvider is a struct that provides all methods required by wallet command.
type WalletProvider struct {
	client coreapi.CoreClient
}

// NewWalletProvider creates new WalletProvider.
func NewWalletProvider(client coreapi.CoreClient) *WalletProvider {
	return &WalletProvider{client: client}
}

func (p *WalletProvider) List() ([]common.Address, error) { return nil, nil }
func (p *WalletProvider) Create(passphrase string) (common.Address, error) {
	return common.Address{}, nil
}
func (p *WalletProvider) Delete(address common.Address, passphrase string) error { return nil }
func (p *WalletProvider) Export(address common.Address, passphrase string) ([]byte, error) {
	return nil, nil
}
func (p *WalletProvider) Import(address common.Address, passphrase string, account []byte) error {
	return nil
}
func (p *WalletProvider) Sign(address common.Address, passphrase string) (*types.Transaction, error) {
	return nil, nil
}
