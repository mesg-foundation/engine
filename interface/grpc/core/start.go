package core

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/core"
)

// StartService starts a service.
func (s *Server) StartService(ctx context.Context, request *coreapi.StartServiceRequest) (*coreapi.StartServiceReply, error) {
	return &coreapi.StartServiceReply{}, s.api.StartService(request.ServiceID)
}
