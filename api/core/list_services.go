package core

import (
	"github.com/mesg-foundation/core/database/services"
	"golang.org/x/net/context"
)

// ListServices return all services from the database
func (s *Server) ListServices(ctx context.Context, request *ListServicesRequest) (reply *ListServicesReply, err error) {
	services, err := services.All()
	if err != nil {
		return
	}
	reply = &ListServicesReply{
		Services: services,
	}
	return
}
