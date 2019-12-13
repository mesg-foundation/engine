package api

import (
	"context"

	"github.com/mesg-foundation/engine/account"
	protobuf_api "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/sdk"
)

// AccountServer is the type to aggregate all Account APIs.
type AccountServer struct {
	sdk *sdk.SDK
}

// NewAccountServer creates a new ServiceServer.
func NewAccountServer(sdk *sdk.SDK) *AccountServer {
	return &AccountServer{sdk: sdk}
}

// List accounts.
func (s *AccountServer) List(ctx context.Context, request *protobuf_api.ListAccountRequest) (*protobuf_api.ListAccountResponse, error) {
	accounts, err := s.sdk.Account.List()
	if err != nil {
		return nil, err
	}
	return &protobuf_api.ListAccountResponse{Accounts: accounts}, nil
}

// Create creates a new account from service.
func (s *AccountServer) Create(ctx context.Context, request *protobuf_api.CreateAccountRequest) (*protobuf_api.CreateAccountResponse, error) {
	account, mnemonic, err := s.sdk.Account.Create(request.Name, request.Password)
	if err != nil {
		return nil, err
	}
	return &protobuf_api.CreateAccountResponse{
		Address:  account.Address,
		Mnemonic: mnemonic,
	}, nil
}

// Get retrives account.
func (s *AccountServer) Get(ctx context.Context, request *protobuf_api.GetAccountRequest) (*account.Account, error) {
	return s.sdk.Account.Get(request.Name)
}

// Delete an account
func (s *AccountServer) Delete(ctx context.Context, request *protobuf_api.DeleteAccountRequest) (*protobuf_api.DeleteAccountResponse, error) {
	credname, credpassword, err := GetCredentialFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.sdk.Account.Delete(credname, credpassword); err != nil {
		return nil, err
	}
	return &protobuf_api.DeleteAccountResponse{}, nil
}
