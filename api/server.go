package api

import (
	"errors"
	"log"
	"net"
	"os"

	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/api/service"
	"github.com/mesg-foundation/core/mesg"
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

	log.Println("Server listens on", s.listener.Addr())

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
	m, err := mesg.New()
	if err != nil {
		log.Fatal(err)
	}

	coreServer, err := core.NewServer(core.MESGOption(m))
	if err != nil {
		log.Fatal(err)
	}

	service.RegisterServiceServer(s.instance, &service.Server{})
	core.RegisterCoreServer(s.instance, coreServer)

	reflection.Register(s.instance)
}
