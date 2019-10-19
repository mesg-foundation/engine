package api

import (
	"context"

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
func (s *InstanceServer) List(ctx context.Context, request *protobuf_api.InstanceServiceListRequest) (*protobuf_api.InstanceServiceListResponse, error) {
	instances, err := s.sdk.Instance.List(&instancesdk.Filter{ServiceHash: request.ServiceHash})
	if err != nil {
		return nil, err
	}
	return &protobuf_api.InstanceServiceListResponse{Instances: instances}, nil
}

// Create creates a new instance from service.
func (s *InstanceServer) Create(ctx context.Context, request *protobuf_api.InstanceServiceCreateRequest) (*protobuf_api.InstanceServiceCreateResponse, error) {
	i, err := s.sdk.Instance.Create(request.ServiceHash, request.Env)
	if err != nil {
		return nil, err
	}
	return &protobuf_api.InstanceServiceCreateResponse{Hash: i.Hash}, nil
}

// Get retrives instance.
func (s *InstanceServer) Get(ctx context.Context, request *protobuf_api.InstanceServiceGetRequest) (*protobuf_api.InstanceServiceGetResponse, error) {
	instance, err := s.sdk.Instance.Get(request.Hash)
	if err != nil {
		return nil, err
	}
	return &protobuf_api.InstanceServiceGetResponse{Instance: instance}, nil
}

// Delete an instance
func (s *InstanceServer) Delete(ctx context.Context, request *protobuf_api.InstanceServiceDeleteRequest) (*protobuf_api.InstanceServiceDeleteResponse, error) {
	if err := s.sdk.Instance.Delete(request.Hash, request.DeleteData); err != nil {
		return nil, err
	}
	return &protobuf_api.InstanceServiceDeleteResponse{}, nil
}
