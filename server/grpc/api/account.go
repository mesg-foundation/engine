package api

import (
	"context"

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
func (s *AccountServer) List(ctx context.Context, request *protobuf_api.AccountServiceListRequest) (*protobuf_api.AccountServiceListResponse, error) {
	accounts, err := s.sdk.Account.List()
	if err != nil {
		return nil, err
	}
	return &protobuf_api.AccountServiceListResponse{Accounts: accounts}, nil
}

// Create creates a new account from service.
func (s *AccountServer) Create(ctx context.Context, request *protobuf_api.AccountServiceCreateRequest) (*protobuf_api.AccountServiceCreateResponse, error) {
	account, mnemonic, err := s.sdk.Account.Create(request.Name, request.Password)
	if err != nil {
		return nil, err
	}
	return &protobuf_api.AccountServiceCreateResponse{
		Address:  account.Address,
		Mnemonic: mnemonic,
	}, nil
}

// Get retrives account.
func (s *AccountServer) Get(ctx context.Context, request *protobuf_api.AccountServiceGetRequest) (*protobuf_api.AccountServiceGetResponse, error) {
	account, err := s.sdk.Account.Get(request.Name)
	if err != nil {
		return nil, err
	}
	return &protobuf_api.AccountServiceGetResponse{Account: account}, nil
}

// Delete an account
func (s *AccountServer) Delete(ctx context.Context, request *protobuf_api.AccountServiceDeleteRequest) (*protobuf_api.AccountServiceDeleteResponse, error) {
	credname, credpassword, err := GetCredentialFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.sdk.Account.Delete(credname, credpassword); err != nil {
		return nil, err
	}
	return &protobuf_api.AccountServiceDeleteResponse{}, nil
}
