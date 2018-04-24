package service

import (
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/types"
	"golang.org/x/net/context"
)

// Stop a service
func (s *Server) Stop(ctx context.Context, request *types.StopServiceRequest) (reply *types.ServiceReply, err error) {
	service := service.New(request.Service)
	err = service.Stop()
	reply = &types.ServiceReply{}
	return
}
