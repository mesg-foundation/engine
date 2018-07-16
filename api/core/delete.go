package core

import (
	"github.com/mesg-foundation/core/database/services"
	"context"
)

// DeleteService delete a service in the database and eventually stop the docker of this service
func (s *Server) DeleteService(ctx context.Context, request *DeleteServiceRequest) (reply *DeleteServiceReply, err error) {
	service, err := services.Get(request.ServiceID)
	if err != nil {
		return
	}
	err = service.Stop()
	if err != nil {
		return
	}
	err = services.Delete(request.ServiceID)
	if err != nil {
		return
	}
	reply = &DeleteServiceReply{}
	return
}
