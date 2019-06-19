package grpc

import (
	"context"
	"errors"

	protobuf_api "github.com/mesg-foundation/core/protobuf/api"
	"github.com/mesg-foundation/core/sdk"
	"github.com/mesg-foundation/core/server/grpc/api"
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
	srv := api.FromProtoService(request.Definition)
	if err := s.sdk.Service.Create(srv); err != nil {
		return nil, err
	}
	return &protobuf_api.CreateServiceResponse{Sid: srv.Sid, Hash: srv.Hash}, nil
}

// Delete deletes service by hash or sid.
func (s *ServiceServer) Delete(ctx context.Context, request *protobuf_api.DeleteServiceRequest) (*protobuf_api.DeleteServiceResponse, error) {
	srv, err := s.sdk.GetService(request.HashOrSid)
	if err != nil {
		return nil, err
	}
	// first, check if service has any running instances.
	instances, err := s.sdk.Instance.GetAllByService(srv.Hash)
	if err != nil {
		return nil, err
	}
	if len(instances) > 0 {
		return nil, errors.New("service has running instances. in order to delete the service, stop its instances first")
	}
	return &protobuf_api.DeleteServiceResponse{}, s.sdk.Service.Delete(request.HashOrSid)
}
