package service

import (
	"log"

	"golang.org/x/net/context"
)

func (s *Server) Start(ctx context.Context, request *StartRequest) (reply *StartReply, err error) {
	log.Println("receive start", request)
	reply = &StartReply{}
	return
}
