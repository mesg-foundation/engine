package core

import (
	"context"

	"github.com/mesg-foundation/core/database/services"
)

// GetService fetches a service from the database and returns it.
func (s *Server) GetService(ctx context.Context, request *GetServiceRequest) (*GetServiceReply, error) {
	service, err := services.Get(request.ServiceID)
	return &GetServiceReply{
		Service: &service,
	}, err
}
