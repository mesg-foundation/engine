package api

import (
	"context"

	"github.com/mesg-foundation/engine/instance"
	protobuf_api "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/sdk"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
)

// InstanceServer is the type to aggregate all Instance APIs.
type InstanceServer struct {
	sdk *sdk.SDK
}

// NewInstanceServer creates a new ServiceServer.
func NewInstanceServer(sdk *sdk.SDK) *InstanceServer {
	return &InstanceServer{sdk: sdk}
}

// List instances.
func (s *InstanceServer) List(ctx context.Context, request *protobuf_api.ListInstanceRequest) (*protobuf_api.ListInstanceResponse, error) {
	instances, err := s.sdk.Instance.List(&instancesdk.Filter{ServiceHash: request.ServiceHash})
	if err != nil {
		return nil, err
	}
	return &protobuf_api.ListInstanceResponse{Instances: instances}, nil
}

// Get retrives instance.
func (s *InstanceServer) Get(ctx context.Context, request *protobuf_api.GetInstanceRequest) (*instance.Instance, error) {
	return s.sdk.Instance.Get(request.Hash)
}
