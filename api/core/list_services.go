package core

import (
	"context"

	"github.com/mesg-foundation/core/database/services"
)

// ListServices return all services from the database
func (s *Server) ListServices(ctx context.Context, request *ListServicesRequest) (*ListServicesReply, error) {
	services, err := services.All()
	if err != nil {
		return nil, err
	}
	return &ListServicesReply{
		Services: services,
	}, nil
}
