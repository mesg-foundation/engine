package grpc

import (
	"errors"
	"net"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/interface/grpc/core"
	"github.com/mesg-foundation/core/interface/grpc/service"
	"github.com/sirupsen/logrus"
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
	s.instance = grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_logrus.StreamServerInterceptor(logrus.NewEntry(logrus.StandardLogger())),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrus.StandardLogger())),
		)),
	)
	if err := s.register(); err != nil {
		return err
	}

	logrus.Info("Server listens on ", s.listener.Addr())

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
func (s *Server) register() error {
	a, err := api.New()
	if err != nil {
		return err
	}

	coreServer, err := core.NewServer(core.APIOption(a))
	if err != nil {
		return err
	}

	serviceServer, err := service.NewServer(service.APIOption(a))
	if err != nil {
		return err
	}

	service.RegisterServiceServer(s.instance, serviceServer)
	core.RegisterCoreServer(s.instance, coreServer)

	reflection.Register(s.instance)
	return nil
}
