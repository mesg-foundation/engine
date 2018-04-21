package service

import (
	"log"

	types "github.com/mesg-foundation/application/types"
	"golang.org/x/net/context"
)

func (s *Server) Start(ctx context.Context, request *types.StartServiceRequest) (reply *types.ServiceReply, err error) {
	log.Println("receive emit", request.Service)
	reply = &types.ServiceReply{}
	return
}
