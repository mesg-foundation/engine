package core

import (
	"context"
)

// DeleteService stops and deletes service serviceID.
func (s *Server) DeleteService(ctx context.Context, request *DeleteServiceRequest) (*DeleteServiceReply, error) {
	return &DeleteServiceReply{}, s.api.DeleteService(request.ServiceID)
}
