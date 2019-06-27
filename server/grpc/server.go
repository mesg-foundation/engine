package grpc

import (
	"errors"
	"net"
	"sync"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	protobuf_api "github.com/mesg-foundation/core/protobuf/api"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/sdk"
	"github.com/mesg-foundation/core/server/grpc/api"
	"github.com/mesg-foundation/core/server/grpc/core"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// Server contains the server config.
type Server struct {
	sdk *sdk.SDK

	instance *grpc.Server
	closed   bool
	mi       sync.Mutex // protects startup.

	network string
	address string
}

// New returns a new gRPC server.
func New(address string, sdk *sdk.SDK) *Server {
	return &Server{
		sdk:     sdk,
		address: address,
		network: "tcp",
	}
}

// listen listens for connections.
func (s *Server) listen() (net.Listener, error) {
	s.mi.Lock()
	defer s.mi.Unlock()

	if s.closed {
		return nil, errors.New("already closed")
	}

	ln, err := net.Listen(s.network, s.address)
	if err != nil {
		return nil, err
	}

	// Keep alive prevents Docker network to drop TCP idle connections after 15 minutes.
	// See: https://forum.mesg.com/t/solution-summary-for-docker-dropping-connections-after-15-min/246
	keepaliveOpt := grpc.KeepaliveParams(keepalive.ServerParameters{
		Time: 1 * time.Minute,
	})
	s.instance = grpc.NewServer(
		keepaliveOpt,
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_logrus.StreamServerInterceptor(logrus.NewEntry(logrus.StandardLogger())),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrus.StandardLogger())),
		)),
	)
	s.register()
	logrus.Info("server listens on ", ln.Addr())
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
func (s *Server) register() {
	coreapi.RegisterCoreServer(s.instance, core.NewServer(s.sdk))

	protobuf_api.RegisterEventServer(s.instance, api.NewEventServer(s.sdk))
	protobuf_api.RegisterExecutionServer(s.instance, api.NewExecutionServer(s.sdk))
	protobuf_api.RegisterInstanceServer(s.instance, api.NewInstanceServer(s.sdk))
	protobuf_api.RegisterServiceXServer(s.instance, api.NewServiceServer(s.sdk))

	reflection.Register(s.instance)
}
