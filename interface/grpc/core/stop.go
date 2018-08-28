package core

import (
	"context"
)

// StopService stops a service.
func (s *Server) StopService(ctx context.Context, request *StopServiceRequest) (*StopServiceReply, error) {
	return &StopServiceReply{}, s.api.StopService(request.ServiceID)
}
