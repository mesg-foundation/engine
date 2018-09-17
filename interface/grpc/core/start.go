package core

import (
	"context"

	"github.com/mesg-foundation/core/protobuf/core"
)

// StartService starts a service.
func (s *Server) StartService(ctx context.Context, request *core.StartServiceRequest) (*core.StartServiceReply, error) {
	return &core.StartServiceReply{}, s.api.StartService(request.ServiceID)
}
