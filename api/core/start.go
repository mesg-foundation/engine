package core

import (
	"context"

	"github.com/mesg-foundation/core/database/services"
)

// StartService fetches a service from the database and starts it.
func (s *Server) StartService(ctx context.Context, request *StartServiceRequest) (*StartServiceReply, error) {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return nil, err
	}
	if _, err = service.Start(); err != nil {
		return nil, err
	}
	return &StartServiceReply{}, nil
}
