package core

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/core"
)

// GetService returns service serviceID.
func (s *Server) GetService(ctx context.Context, request *core.GetServiceRequest) (*core.GetServiceReply, error) {
	srv, err := s.api.GetService(request.ServiceID)
	if err != nil {
		return nil, err
	}
	return &core.GetServiceReply{Service: toProtoService(srv)}, nil
}
