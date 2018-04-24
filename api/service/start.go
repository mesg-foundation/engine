package service

import (
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/types"
	"golang.org/x/net/context"
)

// Start a service
func (s *Server) Start(ctx context.Context, request *types.StartServiceRequest) (reply *types.ServiceReply, err error) {
	service := service.New(request.Service)
	_, err = service.Start()
	reply = &types.ServiceReply{}
	return
}
