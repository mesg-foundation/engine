package grpc

import (
	"context"

	protobuf_api "github.com/mesg-foundation/core/protobuf/api"
	"github.com/mesg-foundation/core/protobuf/definition"
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

// Get returns service from given hash.
func (s *ServiceServer) Get(ctx context.Context, req *protobuf_api.GetServiceRequest) (*definition.Service, error) {
	exec, err := s.sdk.Service.Get(req.Hash)
	if err != nil {
		return nil, err
	}
	return api.ToProtoService(exec), nil
}
