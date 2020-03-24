package api

import (
	"context"
	"fmt"

	"github.com/mesg-foundation/engine/cosmos"
	protobuf_api "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/service"
)

// ServiceServer is the type to aggregate all Service APIs.
type ServiceServer struct {
	mc *cosmos.ModuleClient
}

// NewServiceServer creates a new ServiceServer.
func NewServiceServer(mc *cosmos.ModuleClient) *ServiceServer {
	return &ServiceServer{mc: mc}
}

// Create creates a new service from definition.
func (s *ServiceServer) Create(ctx context.Context, req *protobuf_api.CreateServiceRequest) (*protobuf_api.CreateServiceResponse, error) {
	srv, err := s.mc.CreateService(req)
	if err != nil {
		return nil, err
	}
	return &protobuf_api.CreateServiceResponse{Hash: srv.Hash}, nil
}

// Get returns service from given hash.
func (s *ServiceServer) Get(ctx context.Context, req *protobuf_api.GetServiceRequest) (*service.Service, error) {
	return s.mc.GetService(req.Hash)
}

// List returns all services.
func (s *ServiceServer) List(ctx context.Context, req *protobuf_api.ListServiceRequest) (*protobuf_api.ListServiceResponse, error) {
	services, err := s.mc.ListService()
	if err != nil {
		return nil, err
	}

	return &protobuf_api.ListServiceResponse{Services: services}, nil
}

// Exists returns if a service already exists.
func (s *ServiceServer) Exists(ctx context.Context, req *protobuf_api.ExistsServiceRequest) (*protobuf_api.ExistsServiceResponse, error) {
	exist, err := s.mc.ExistService(req.Hash)
	if err != nil {
		return nil, err
	}
	return &protobuf_api.ExistsServiceResponse{Exists: exist}, nil
}

// Hash returns the calculated hash of a service request.
func (s *ServiceServer) Hash(ctx context.Context, req *protobuf_api.CreateServiceRequest) (*protobuf_api.HashServiceResponse, error) {
	return nil, fmt.Errorf("not implemented anymore, use LCD")
}
