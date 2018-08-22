package core

import (
	"context"
)

// GetService returns service serviceID.
func (s *Server) GetService(ctx context.Context, request *GetServiceRequest) (*GetServiceReply, error) {
	service, err := s.api.GetService(request.ServiceID)
	return &GetServiceReply{Service: service}, err
}
