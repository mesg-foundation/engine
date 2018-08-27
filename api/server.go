package api

import (
	"net"
	"os"
	"sync"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/api/service"
	"github.com/mesg-foundation/core/mesg"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server contains the server config.
type Server struct {
	instance *grpc.Server
	closed   bool
	mi       sync.Mutex // protects startup.

	Network string
	Address string
}

// listen listens for connections.
func (s *Server) listen() (net.Listener, error) {
	s.mi.Lock()
	defer s.mi.Unlock()

	if s.closed {
		return nil, alreadyClosedError{}
	}

	if s.Network == "unix" {
		os.Remove(s.Address)
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
	m, err := mesg.New()
	if err != nil {
		return err
	}

	coreServer, err := core.NewServer(core.MESGOption(m))
	if err != nil {
		return err
	}

	service.RegisterServiceServer(s.instance, &service.Server{})
	core.RegisterCoreServer(s.instance, coreServer)

	reflection.Register(s.instance)
	return nil
}

type alreadyClosedError struct{}

func (e alreadyClosedError) Error() string {
	return "already closed"
}
