package core

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/core"
)

// ListServices lists services.
func (s *Server) ListServices(ctx context.Context, request *core.ListServicesRequest) (*core.ListServicesReply, error) {
	services, err := s.api.ListServices()
	if err != nil {
		return nil, err
	}
	return &core.ListServicesReply{Services: toProtoServices(services)}, nil
}
