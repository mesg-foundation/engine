package api

import (
	"net"

	"github.com/mesg-foundation/application/api/service"
	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

var server *grpc.Server

// StartServer starts the server
func StartServer() (err error) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return
	}
	server = grpc.NewServer()
	register(server)
	// Register reflection service on gRPC server.
	// reflection.Register(s)
	err = server.Serve(listener)
	if err != nil {
		return
	}
	return
}

// StopServer stops the server (if exist)
func StopServer() {
	if server != nil {
		server.Stop()
		server = nil
	}
}

func register(server *grpc.Server) {
	apiService.Register(server)
}
