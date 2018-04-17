package api

import (
	"net"

	"github.com/mesg-foundation/application/api/service"
	"github.com/mesg-foundation/application/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server is the main struct that contain the server config
type Server struct {
	instance *grpc.Server
	Network  string
	Address  string
}

// network returns the Server's network or a default
func (s *Server) network() (network string) {
	network = s.Network
	if network == "" {
		network = config.Api.Server.Network()
	}
	return
}

// address returns the Server's address or a default
func (s *Server) address() (address string) {
	address = s.Address
	if address == "" {
		address = config.Api.Server.Address()
	}
	return
}

// Start starts the server
func (s *Server) Start() (err error) {
	listener, err := net.Listen(s.network(), s.address())
	if err != nil {
		return
	}
	s.instance = grpc.NewServer()
	s.register()
	reflection.Register(s.instance)
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
	service.RegisterServiceServer(s.instance, &service.Server{})
}
