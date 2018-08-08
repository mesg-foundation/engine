package api

import (
	"errors"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/api/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server contains the server config.
type Server struct {
	instance *grpc.Server
	listener net.Listener
	Network  string
	Address  string
}

// Serve starts the server and listens for client connections.
func (s *Server) Serve() error {
	if s.listener != nil {
		return errors.New("Server already running")
	}

	if s.Network == "unix" {
		os.Remove(s.Address)
	}
	listener, err := net.Listen(s.Network, s.Address)
	if err != nil {
		return err
	}

	s.listener = listener
	s.instance = grpc.NewServer()
	s.register()

	log.Info("Server listens on", s.listener.Addr())

	// TODO: check if server still on after a connection throw an error. otherwise, add a for around serve
	return s.instance.Serve(s.listener)
}

// Stop stops the server.
func (s *Server) Stop() {
	if s.instance != nil {
		s.instance.Stop()
		s.instance = nil
	}
}

// register all server
func (s *Server) register() {
	service.RegisterServiceServer(s.instance, &service.Server{})
	core.RegisterCoreServer(s.instance, &core.Server{})

	reflection.Register(s.instance)
}
