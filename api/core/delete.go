package core

import (
	"github.com/mesg-foundation/core/database/services"
	"golang.org/x/net/context"
)

// DeleteService delete a service in the database
func (s *Server) DeleteService(ctx context.Context, request *DeleteServiceRequest) (reply *DeleteServiceReply, err error) {
	err = services.Delete(request.ServiceID)
	if err != nil {
		return
	}
	reply = &DeleteServiceReply{}
	return
}
