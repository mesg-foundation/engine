package core

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/coreapi"
)

// GetService returns service serviceID.
func (s *Server) GetService(ctx context.Context, request *coreapi.GetServiceRequest) (*coreapi.GetServiceReply, error) {
	srv, err := s.api.GetService(request.ServiceID)
	if err != nil {
		return nil, err
	}
	return &coreapi.GetServiceReply{Service: toProtoService(srv)}, nil
}
