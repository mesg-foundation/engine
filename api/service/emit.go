package service

import (
	"log"

	"golang.org/x/net/context"
)

func (s *Server) Emit(ctx context.Context, request *EmitRequest) (reply *EmitReply, err error) {
	log.Println("receive emit", request)
	reply = &EmitReply{}
	return
}
