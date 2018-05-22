package core

import (
	"golang.org/x/net/context"
)

// Start a service
func (s *Server) StartService(ctx context.Context, request *StartServiceRequest) (reply *StartServiceReply, err error) {
	service := request.Service
	_, err = service.Start()
	reply = &StartServiceReply{}
	return
}
