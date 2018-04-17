package service

import (
	"log"

	"golang.org/x/net/context"
)

func (s *Server) Stop(ctx context.Context, request *StopRequest) (reply *StopReply, err error) {
	log.Println("receive stop", request)
	reply = &StopReply{}
	return
}
