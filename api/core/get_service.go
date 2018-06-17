package core

import (
	"github.com/mesg-foundation/core/database/services"
	"golang.org/x/net/context"
)

// GetService fetch a service in the db and return ot
func (s *Server) GetService(ctx context.Context, request *GetServiceRequest) (reply *GetServiceReply, err error) {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return
	}
	reply = &GetServiceReply{
		Service: &service,
	}
	return
}
