package api

import (
	"context"

	protobuf_api "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/sdk"
)

// OwnershipServer is the type to aggregate all Ownership APIs.
type OwnershipServer struct {
	sdk *sdk.SDK
}

// NewOwnershipServer creates a new OwnershipServer.
func NewOwnershipServer(sdk *sdk.SDK) *OwnershipServer {
	return &OwnershipServer{sdk: sdk}
}

// List returns all ownerships.
func (s *OwnershipServer) List(ctx context.Context, req *protobuf_api.OwnershipServiceListRequest) (*protobuf_api.OwnershipServiceListResponse, error) {
	ownerships, err := s.sdk.Ownership.List()
	if err != nil {
		return nil, err
	}

	return &protobuf_api.OwnershipServiceListResponse{Ownerships: ownerships}, nil
}
