package api

import (
	"context"
	"fmt"

	"github.com/mesg-foundation/engine/cosmos"
	instancepb "github.com/mesg-foundation/engine/instance"
	protobuf_api "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/sdk"
	"github.com/mesg-foundation/engine/x/instance"
)

// InstanceServer is the type to aggregate all Instance APIs.
type InstanceServer struct {
	sdk    *sdk.SDK
	client *cosmos.Client
}

// NewInstanceServer creates a new ServiceServer.
func NewInstanceServer(sdk *sdk.SDK, client *cosmos.Client) *InstanceServer {
	return &InstanceServer{
		sdk:    sdk,
		client: client,
	}
}

// List instances.
func (s *InstanceServer) List(ctx context.Context, request *protobuf_api.ListInstanceRequest) (*protobuf_api.ListInstanceResponse, error) {
	var instances []*instancepb.Instance
	if err := s.client.QueryJSON(fmt.Sprintf("custom/%s/%s", instance.QuerierRoute, instance.QueryListInstances), request.Filter, &instances); err != nil {
		return nil, err
	}
	return &protobuf_api.ListInstanceResponse{Instances: instances}, nil
}

// Get retrives instance.
func (s *InstanceServer) Get(ctx context.Context, request *protobuf_api.GetInstanceRequest) (*instancepb.Instance, error) {
	var inst instancepb.Instance
	err := s.client.QueryJSON(
		fmt.Sprintf("custom/%s/%s/%s", instance.QuerierRoute, instance.QueryGetInstance, request.Hash),
		nil, &inst)
	if err != nil {
		return nil, err
	}
	return &inst, nil
}
