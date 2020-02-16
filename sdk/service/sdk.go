package servicesdk

import (
	"fmt"

	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/api"
	servicepb "github.com/mesg-foundation/engine/service"
	"github.com/mesg-foundation/engine/x/service"
)

// SDK is the service sdk.
type SDK struct {
	client *cosmos.Client
}

// New returns the service sdk.
func New(client *cosmos.Client) *SDK {
	return &SDK{client: client}
}

// Create creates a new service from definition.
func (s *SDK) Create(req *api.CreateServiceRequest) (*servicepb.Service, error) {
	acc, err := s.client.GetAccount()
	if err != nil {
		return nil, err
	}
	msg := service.NewMsgCreateService(acc.GetAddress(), req)
	tx, err := s.client.BuildAndBroadcastMsg(msg)
	if err != nil {
		return nil, err
	}
	return s.Get(tx.Data)
}

// Get returns the service that matches given hash.
func (s *SDK) Get(hash hash.Hash) (*servicepb.Service, error) {
	var se servicepb.Service
	if err := s.client.QueryJSON(fmt.Sprintf("custom/%s/%s/%s", service.QuerierRoute, service.QueryGetService, hash), nil, &se); err != nil {
		return nil, err
	}
	return &se, nil
}

// List returns all services.
func (s *SDK) List() ([]*servicepb.Service, error) {
	var services []*servicepb.Service
	if err := s.client.QueryJSON(fmt.Sprintf("custom/%s/%s", service.QuerierRoute, service.QueryListService), nil, &services); err != nil {
		return nil, err
	}
	return services, nil
}

// Exists returns if a service already exists.
func (s *SDK) Exists(hash hash.Hash) (bool, error) {
	var exists bool
	if err := s.client.QueryJSON(fmt.Sprintf("custom/%s/%s/%s", service.QuerierRoute, service.QueryExistService, hash), nil, &exists); err != nil {
		return false, err
	}
	return exists, nil
}

// Hash returns the calculate hash of a service.
func (s *SDK) Hash(req *api.CreateServiceRequest) (hash.Hash, error) {
	var h hash.Hash
	if err := s.client.QueryJSON(fmt.Sprintf("custom/%s/%s", service.QuerierRoute, service.QueryHashService), req, &h); err != nil {
		return nil, err
	}
	return h, nil
}
