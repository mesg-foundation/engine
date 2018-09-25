package grpc

import (
	"net"
	"sync"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/interface/grpc/core"
	"github.com/mesg-foundation/core/interface/grpc/service"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/protobuf/serviceapi"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server contains the server config.
type Server struct {
	instance *grpc.Server
	closed   bool
	mi       sync.Mutex // protects startup.

	Network   string
	Address   string
	ServiceDB *database.ServiceDB
}

// listen listens for connections.
func (s *Server) listen() (net.Listener, error) {
	s.mi.Lock()
	defer s.mi.Unlock()

	if s.closed {
		return nil, &alreadyClosedError{}
	}

	ln, err := net.Listen(s.Network, s.Address)
	if err != nil {
		return nil, err
	}

	s.instance = grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_logrus.StreamServerInterceptor(logrus.NewEntry(logrus.StandardLogger())),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrus.StandardLogger())),
		)),
	)
	if err := s.register(); err != nil {
		return nil, err
	}

	logrus.Info("Server listens on ", ln.Addr())
	return ln, nil
}

// Serve starts the server and listens for client connections.
func (s *Server) Serve() error {
	ln, err := s.listen()
	if err != nil {
		return err
	}

	// TODO: check if server still on after a connection throw an error. otherwise, add a for around serve
	return s.instance.Serve(ln)
}

// Close gracefully closes the server.
func (s *Server) Close() {
	s.mi.Lock()
	defer s.mi.Unlock()
	if s.closed {
		return
	}
	if s.instance != nil {
		s.instance.GracefulStop()
	}
	s.closed = true
}

// register all server
func (s *Server) register() error {
	a, err := api.New(api.DatabaseOption(s.ServiceDB))
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

	serviceapi.RegisterServiceServer(s.instance, serviceServer)
	coreapi.RegisterCoreServer(s.instance, coreServer)

	reflection.Register(s.instance)
	return nil
}

type alreadyClosedError struct{}

func (e *alreadyClosedError) Error() string {
	return "already closed"
}
