package api

import (
	"context"
	"errors"

	"github.com/mesg-foundation/core/hash"
	protobuf_api "github.com/mesg-foundation/core/protobuf/api"
	"github.com/mesg-foundation/core/protobuf/definition"
	"github.com/mesg-foundation/core/sdk"
	instancesdk "github.com/mesg-foundation/core/sdk/instance"
)

// ServiceServer is the type to aggregate all Service APIs.
type ServiceServer struct {
	sdk *sdk.SDK
}

// NewServiceServer creates a new ServiceServer.
func NewServiceServer(sdk *sdk.SDK) *ServiceServer {
	return &ServiceServer{sdk: sdk}
}

// Create creates a new service from definition.
func (s *ServiceServer) Create(ctx context.Context, request *protobuf_api.CreateServiceRequest) (*protobuf_api.CreateServiceResponse, error) {
	definition := fromProtoService(request.Definition)
	srv, err := s.sdk.Service.Create(definition)
	if err != nil {
		return nil, err
	}
	return &protobuf_api.CreateServiceResponse{Sid: srv.Sid, Hash: srv.Hash.String()}, nil
}

// Delete deletes service by hash or sid.
func (s *ServiceServer) Delete(ctx context.Context, request *protobuf_api.DeleteServiceRequest) (*protobuf_api.DeleteServiceResponse, error) {
	hash, err := hash.Decode(request.Hash)
	if err != nil {
		return nil, err
	}

	srv, err := s.sdk.Service.Get(hash)
	if err != nil {
		return nil, err
	}
	// first, check if service has any running instances.
	instances, err := s.sdk.Instance.List(&instancesdk.Filter{ServiceHash: srv.Hash})
	if err != nil {
		return nil, err
	}
	if len(instances) > 0 {
		return nil, errors.New("service has running instances. in order to delete the service, stop its instances first")
	}
	return &protobuf_api.DeleteServiceResponse{}, s.sdk.Service.Delete(hash)
}

// Get returns service from given hash.
func (s *ServiceServer) Get(ctx context.Context, req *protobuf_api.GetServiceRequest) (*definition.Service, error) {
	hash, err := hash.Decode(req.Hash)
	if err != nil {
		return nil, err
	}

	service, err := s.sdk.Service.Get(hash)
	if err != nil {
		return nil, err
	}
	return toProtoService(service), nil
}

// List returns all services.
func (s *ServiceServer) List(ctx context.Context, req *protobuf_api.ListServiceRequest) (*protobuf_api.ListServiceResponse, error) {
	services, err := s.sdk.Service.List()
	if err != nil {
		return nil, err
	}

	resp := &protobuf_api.ListServiceResponse{}
	for _, service := range services {
		resp.Services = append(resp.Services, toProtoService(service))
	}

	return resp, nil
}
