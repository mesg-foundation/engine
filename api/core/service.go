package core

import (
	"github.com/mesg-foundation/core/database/services"
	"golang.org/x/net/context"
)

// Service fetch a service in the db and return ot
func (s *Server) Service(ctx context.Context, request *ServiceRequest) (reply *ServiceReply, err error) {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return
	}
	reply = &ServiceReply{
		Service: &service,
	}
	return
}
