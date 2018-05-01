package api

import (
	"log"
	"net"
	"os"

	"github.com/mesg-foundation/core/api/event"
	"github.com/mesg-foundation/core/api/service"
	"github.com/mesg-foundation/core/api/task"
	"github.com/mesg-foundation/core/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server is the main struct that contain the server config
type Server struct {
	instance *grpc.Server
	Network  string
	Address  string
}

// Serve starts the server and listen for client connections
func (s *Server) Serve() (err error) {
	os.Remove(s.Address)
	listener, err := net.Listen(s.Network, s.Address)
	if err != nil {
		return
	}

	s.instance = grpc.NewServer()
	s.register()

	log.Println("Server listens on", listener.Addr())

	// TODO: check if server still on after a connection throw an error. otherwise, add a for around serve
	err = s.instance.Serve(listener)
	if err != nil {
		return
	}
	return
}

// Stop stops the server (if exist)
func (s *Server) Stop() {
	if s.instance != nil {
		s.instance.Stop()
		s.instance = nil
	}
}

// register all server
func (s *Server) register() {

	types.RegisterServiceServer(s.instance, &service.Server{})
	types.RegisterEventServer(s.instance, &event.Server{})
	// s.instance.RegisterService(&service.Descriptions, service.ServerInstance)
	types.RegisterTaskServer(s.instance, &task.Server{})

	reflection.Register(s.instance)
}
