package core

import (
	"context"

	"github.com/mesg-foundation/core/database/services"
)

// GetService fetch a service in the db and return ot
func (s *Server) GetService(ctx context.Context, request *GetServiceRequest) (*GetServiceReply, error) {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return nil, err
	}
	return &GetServiceReply{
		Service: &service,
	}, nil
}
