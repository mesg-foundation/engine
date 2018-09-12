package core

import (
	"context"
)

// GetService returns service serviceID.
func (s *Server) GetService(ctx context.Context, request *GetServiceRequest) (*GetServiceReply, error) {
	srv, err := s.api.GetService(request.ServiceID)
	if err != nil {
		return nil, err
	}
	return &GetServiceReply{Service: toProtoService(srv)}, nil
}
