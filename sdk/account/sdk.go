package servicesdk

import (
	"github.com/mesg-foundation/engine/account"
	"github.com/mesg-foundation/engine/cosmos"
)

// SDK is the account SDK.
type SDK struct {
	kb *cosmos.Keybase
}

// NewSDK returns the account sdk.
func NewSDK(kb *cosmos.Keybase) *SDK {
	sdk := &SDK{
		kb: kb,
	}
	return sdk
}

// Create generates a new mnemonic and its associated account from a name and password.
func (s *SDK) Create(name, password string) (address *account.Account, mnemonic string, err error) {
	// TODO: should throw error if name already exist
	mnemonic, err = s.kb.NewMnemonic()
	if err != nil {
		return nil, "", err
	}
	acct, err := s.kb.CreateAccount(name, mnemonic, "", password, 0, 0)
	if err != nil {
		return nil, "", err
	}
	return &account.Account{
		Name:    acct.GetName(),
		Address: acct.GetAddress().String(),
	}, mnemonic, nil
}

// Delete removes the account corresponding the name and password.
func (s *SDK) Delete(name, password string) error {
	return s.kb.Delete(name, password, false)
}

// Get returns the account from a name.
func (s *SDK) Get(name string) (*account.Account, error) {
	acct, err := s.kb.Get(name)
	if err != nil {
		return nil, err
	}
	return &account.Account{
		Name:    acct.GetName(),
		Address: acct.GetAddress().String(),
	}, nil
}

// List returns all account.
func (s *SDK) List() ([]*account.Account, error) {
	accts, err := s.kb.List()
	if err != nil {
		return nil, err
	}
	accounts := make([]*account.Account, len(accts))
	for i, acct := range accts {
		accounts[i] = &account.Account{
			Name:    acct.GetName(),
			Address: acct.GetAddress().String(),
		}
	}
	return accounts, nil
}
