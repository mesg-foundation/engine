package core

import (
	"context"

	"github.com/mesg-foundation/core/database/services"
)

// ListServices return all services from the database
func (s *Server) ListServices(ctx context.Context, request *ListServicesRequest) (reply *ListServicesReply, err error) {
	services, err := services.All()
	return &ListServicesReply{
		Services: services,
	}, err
}
