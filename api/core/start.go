package core

import (
	"context"
	"github.com/mesg-foundation/core/database/services"
)

// StartService fetch a service in the db and starts it
func (s *Server) StartService(ctx context.Context, request *StartServiceRequest) (reply *StartServiceReply, err error) {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return
	}
	_, err = service.Start()
	reply = &StartServiceReply{}
	return
}
