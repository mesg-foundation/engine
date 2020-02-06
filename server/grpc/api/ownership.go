package api

import (
	"context"

	"github.com/mesg-foundation/engine/cosmos"
	ownershippb "github.com/mesg-foundation/engine/ownership"
	protobuf_api "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/sdk"
	"github.com/mesg-foundation/engine/x/ownership"
)

// OwnershipServer is the type to aggregate all Ownership APIs.
type OwnershipServer struct {
	sdk    *sdk.SDK
	client *cosmos.Client
}

// NewOwnershipServer creates a new OwnershipServer.
func NewOwnershipServer(sdk *sdk.SDK, client *cosmos.Client) *OwnershipServer {
	return &OwnershipServer{
		sdk:    sdk,
		client: client,
	}
}

// List returns all ownerships.
func (s *OwnershipServer) List(ctx context.Context, req *protobuf_api.ListOwnershipRequest) (*protobuf_api.ListOwnershipResponse, error) {
	var ownerships []*ownershippb.Ownership
	if err := s.client.QueryJSON("custom/"+ownership.ModuleName+"/"+ownership.QueryListOwnerships, nil, &ownerships); err != nil {
		return nil, err
	}
	return &protobuf_api.ListOwnershipResponse{Ownerships: ownerships}, nil
}
