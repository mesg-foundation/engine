package apiService

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Emit(ctx context.Context, request *EmitRequest) (reply *EmitReply, err error) {
	log.Println("receive emit", request)
	reply = &EmitReply{Success: true}
	return
}

func (s *server) Subscribe(*SubscribeRequest, Service_SubscribeServer) (err error) {
	return
}

// Register the API for Service
func Register(s *grpc.Server) {
	RegisterServiceServer(s, &server{})
}
