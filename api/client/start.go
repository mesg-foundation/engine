package client

import (
	"golang.org/x/net/context"
)

// Start a service
func (s *Server) StartService(ctx context.Context, request *ServiceRequest) (reply *ErrorReply, err error) {
	service := request.Service
	_, err = service.Start()
	reply = &ErrorReply{}
	return
}
