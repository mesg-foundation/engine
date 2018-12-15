package core

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/coreapi"
)

// StartService starts a service.
func (s *Server) StartService(ctx context.Context, request *coreapi.StartServiceRequest) (*coreapi.StartServiceReply, error) {
	if err := s.requireStake(); err != nil {
		return nil, err
	}
	return &coreapi.StartServiceReply{}, s.api.StartService(request.ServiceID)
}
