package api

import (
	"context"

	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/protobuf/api"
)

// OwnershipServer is the type to aggregate all Ownership APIs.
type OwnershipServer struct {
	mc *cosmos.ModuleClient
}

// NewOwnershipServer creates a new OwnershipServer.
func NewOwnershipServer(mc *cosmos.ModuleClient) *OwnershipServer {
	return &OwnershipServer{mc: mc}
}

// List returns all ownerships.
func (s *OwnershipServer) List(ctx context.Context, req *api.ListOwnershipRequest) (*api.ListOwnershipResponse, error) {
	out, err := s.mc.ListOwnership()
	if err != nil {
		return nil, err
	}
	return &api.ListOwnershipResponse{Ownerships: out}, nil
}
