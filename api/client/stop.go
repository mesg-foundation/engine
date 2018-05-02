package client

import (
	"golang.org/x/net/context"
)

// Stop a service
func (s *Server) StopService(ctx context.Context, request *ServiceRequest) (reply *ErrorReply, err error) {
	service := request.Service
	err = service.Stop()
	reply = &ErrorReply{}
	return
}
