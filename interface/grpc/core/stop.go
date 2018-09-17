package core

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/core"
)

// StopService stops a service.
func (s *Server) StopService(ctx context.Context, request *coreapi.StopServiceRequest) (*coreapi.StopServiceReply, error) {
	return &coreapi.StopServiceReply{}, s.api.StopService(request.ServiceID)
}
