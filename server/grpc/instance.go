package grpc

import (
	"context"

	protobuf_api "github.com/mesg-foundation/core/protobuf/api"
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
	return &protobuf_api.CreateInstanceResponse{
		Hash:        i.Hash,
		ServiceHash: i.ServiceHash,
	}, nil
}
