package core

import (
	"context"

	"github.com/mesg-foundation/core/database/services"
)

// StopService fetches a service from the database and stops it.
func (s *Server) StopService(ctx context.Context, request *StopServiceRequest) (*StopServiceReply, error) {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return nil, err
	}
	return &StopServiceReply{}, service.Stop()
}
