package api

import (
	"context"

	"github.com/mesg-foundation/engine/cosmos"
	instancepb "github.com/mesg-foundation/engine/instance"
	protobuf_api "github.com/mesg-foundation/engine/protobuf/api"
)

// InstanceServer is the type to aggregate all Instance APIs.
type InstanceServer struct {
	mc *cosmos.ModuleClient
}

// NewInstanceServer creates a new ServiceServer.
func NewInstanceServer(mc *cosmos.ModuleClient) *InstanceServer {
	return &InstanceServer{mc: mc}
}

// List instances.
func (s *InstanceServer) List(ctx context.Context, req *protobuf_api.ListInstanceRequest) (*protobuf_api.ListInstanceResponse, error) {
	out, err := s.mc.ListInstance(req)
	if err != nil {
		return nil, err
	}
	return &protobuf_api.ListInstanceResponse{Instances: out}, nil
}

// Get retrives instance.
func (s *InstanceServer) Get(ctx context.Context, req *protobuf_api.GetInstanceRequest) (*instancepb.Instance, error) {
	return s.mc.GetInstance(req.Hash)
}
