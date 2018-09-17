package core

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/core"
)

// StopService stops a service.
func (s *Server) StopService(ctx context.Context, request *core.StopServiceRequest) (*core.StopServiceReply, error) {
	return &core.StopServiceReply{}, s.api.StopService(request.ServiceID)
}
