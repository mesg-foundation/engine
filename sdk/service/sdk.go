package servicesdk

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/service"
)

// SDK is the service sdk.
type SDK struct {
	cdc    *codec.Codec
	client *cosmos.Client
}

// NewSDK returns the service sdk.
func NewSDK(cdc *codec.Codec, client *cosmos.Client) Service {
	sdk := &SDK{
		cdc:    cdc,
		client: client,
	}
	return sdk
}

// Create creates a new service from definition.
func (s *SDK) Create(req *api.CreateServiceRequest) (*service.Service, error) {
	// TODO: pass account by parameters
	accountName, accountPassword := "bob", "pass"
	accNumber, accSeq := uint64(0), uint64(0)
	msg := newMsgCreateService(s.cdc, req)
	tx, err := s.client.BuildAndBroadcastMsg(msg, accountName, accountPassword, accNumber, accSeq)
	if err != nil {
		return nil, err
	}
	return s.Get(tx.Data)
}

// Delete deletes the service by hash.
func (s *SDK) Delete(hash hash.Hash) error {
	// TODO: pass account by parameters
	accountName, accountPassword := "bob", "pass"
	accNumber, accSeq := uint64(0), uint64(0)
	msg := newMsgDeleteService(s.cdc, hash)
	_, err := s.client.BuildAndBroadcastMsg(msg, accountName, accountPassword, accNumber, accSeq)
	return err
}

// Get returns the service that matches given hash.
func (s *SDK) Get(hash hash.Hash) (*service.Service, error) {
	var service service.Service
	if err := s.client.Query("custom/"+backendName+"/get/"+hash.String(), &service); err != nil {
		return nil, err
	}
	return &service, nil
}

// List returns all services.
func (s *SDK) List() ([]*service.Service, error) {
	var services []*service.Service
	if err := s.client.Query("custom/"+backendName+"/list", &services); err != nil {
		return nil, err
	}
	return services, nil
}
