package core

import (
	"golang.org/x/net/context"
)

// Stop a service
func (s *Server) StopService(ctx context.Context, request *StopServiceRequest) (reply *StopServiceReply, err error) {
	service := request.Service
	err = service.Stop()
	reply = &StopServiceReply{}
	return
}
