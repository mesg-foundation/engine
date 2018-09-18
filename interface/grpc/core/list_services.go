package core

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/coreapi"
)

// ListServices lists services.
func (s *Server) ListServices(ctx context.Context, request *coreapi.ListServicesRequest) (*coreapi.ListServicesReply, error) {
	services, err := s.api.ListServices()
	if err != nil {
		return nil, err
	}
	return &coreapi.ListServicesReply{Services: toProtoServices(services)}, nil
}
