package grpc

import (
	"context"

	"github.com/mesg-foundation/core/instance"
	protobuf_api "github.com/mesg-foundation/core/protobuf/api"
	"github.com/mesg-foundation/core/protobuf/definition"
	"github.com/mesg-foundation/core/sdk"
)

// InstanceServer is the type to aggregate all Instance APIs.
type InstanceServer struct {
	sdk *sdk.SDK
}

// NewInstanceServer creates a new ServiceServer.
func NewInstanceServer(sdk *sdk.SDK) *InstanceServer {
	return &InstanceServer{sdk: sdk}
}

// Create creates a new instance from service.
func (s *InstanceServer) Create(ctx context.Context, request *protobuf_api.CreateInstanceRequest) (*protobuf_api.CreateInstanceResponse, error) {
	i, err := s.sdk.Instance.Create(request.ServiceHash, request.Env)
	if err != nil {
		return nil, err
	}
	return &protobuf_api.CreateInstanceResponse{Instance: toProtoInstance(i)}, nil
}

// Get retrives instance.
func (s *InstanceServer) Get(ctx context.Context, request *protobuf_api.GetInstanceRequest) (*protobuf_api.GetInstanceResponse, error) {
	i, err := s.sdk.Instance.Get(request.Hash)
	if err != nil {
		return nil, err
	}
	return &protobuf_api.GetInstanceResponse{Instance: toProtoInstance(i)}, nil
}

// Delete an instance
func (s *InstanceServer) Delete(ctx context.Context, request *protobuf_api.DeleteInstanceRequest) (*protobuf_api.DeleteInstanceResponse, error) {
	err := s.sdk.Instance.Delete(request.Hash)
	if err != nil {
		return nil, err
	}
	return &protobuf_api.DeleteInstanceResponse{
		Hash: request.Hash,
	}, nil
}

func toProtoInstance(i *instance.Instance) *definition.Instance {
	return &definition.Instance{
		Hash:        i.Hash,
		ServiceHash: i.ServiceHash,
	}
}
