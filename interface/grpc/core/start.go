package core

import (
	"context"
)

// StartService starts a service.
func (s *Server) StartService(ctx context.Context, request *StartServiceRequest) (*StartServiceReply, error) {
	return &StartServiceReply{}, s.api.StartService(request.ServiceID)
}
