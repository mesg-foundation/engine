package core

import (
	"github.com/mesg-foundation/core/database/services"
	"golang.org/x/net/context"
)

// StopService fetch a service in db and stop it
func (s *Server) StopService(ctx context.Context, request *StopServiceRequest) (reply *StopServiceReply, err error) {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return
	}
	err = service.Stop()
	reply = &StopServiceReply{}
	return
}
