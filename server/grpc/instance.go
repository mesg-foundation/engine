package grpc

import (
	"context"

	"github.com/mesg-foundation/core/instance"
	protobuf_api "github.com/mesg-foundation/core/protobuf/api"
	"github.com/mesg-foundation/core/protobuf/definition"
	"github.com/mesg-foundation/core/sdk"
	instancesdk "github.com/mesg-foundation/core/sdk/instance"
	"github.com/mr-tron/base58"
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
func (s *InstanceServer) List(ctx context.Context, request *protobuf_api.ListInstancesRequest) (*protobuf_api.ListInstancesResponse, error) {
	instances, err := s.sdk.Instance.List(&instancesdk.Filter{ServiceHash: request.ServiceHash})
	if err != nil {
		return nil, err
	}
	return &protobuf_api.ListInstancesResponse{Instances: toProtoInstances(instances)}, nil
}

// Create creates a new instance from service.
func (s *InstanceServer) Create(ctx context.Context, request *protobuf_api.CreateInstanceRequest) (*protobuf_api.CreateInstanceResponse, error) {
	hash, err := base58.Decode(request.ServiceHash)
	if err != nil {
		return nil, err
	}

	i, err := s.sdk.Instance.Create(hash, request.Env)
	if err != nil {
		return nil, err
	}
	return &protobuf_api.CreateInstanceResponse{
		Hash:        base58.Encode(i.Hash),
		ServiceHash: base58.Encode(i.ServiceHash),
	}, nil
}

// Delete an instance
func (s *InstanceServer) Delete(ctx context.Context, request *protobuf_api.DeleteInstanceRequest) (*protobuf_api.DeleteInstanceResponse, error) {
	hash, err := base58.Decode(request.Hash)
	if err != nil {
		return nil, err
	}
	if err := s.sdk.Instance.Delete(hash); err != nil {
		return nil, err
	}
	return &protobuf_api.DeleteInstanceResponse{
		Hash: request.Hash,
	}, nil
}

func toProtoInstances(instances []*instance.Instance) []*definition.Instance {
	inst := make([]*definition.Instance, len(instances))
	for i, instance := range instances {
		inst[i] = toProtoInstance(instance)
	}
	return inst
}

func toProtoInstance(i *instance.Instance) *definition.Instance {
	return &definition.Instance{
		Hash:        i.Hash,
		ServiceHash: i.ServiceHash,
	}
}
