package core

import (
	"context"

	"github.com/mesg-foundation/core/database/services"
)

// ListServices returns all services from the database.
func (s *Server) ListServices(ctx context.Context, request *ListServicesRequest) (*ListServicesReply, error) {
	services, err := services.All()
	return &ListServicesReply{
		Services: services,
	}, err
}
