package core

import (
	"context"

	"github.com/mesg-foundation/core/database/services"
)

// DeleteService delete a service in the database and eventually stop the docker of this service
func (s *Server) DeleteService(ctx context.Context, request *DeleteServiceRequest) (*DeleteServiceReply, error) {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return nil, err
	}
	if err = service.Stop(); err != nil {
		return nil, err
	}
	if err := services.Delete(request.ServiceID); err != nil {
		return nil, err
	}
	return &DeleteServiceReply{}, nil
}
