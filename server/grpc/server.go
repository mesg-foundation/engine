package grpc

import (
	"net"
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
	instance *grpc.Server
	sdk      *sdk.SDK
}

// New returns a new gRPC server.
func New(sdk *sdk.SDK) *Server {
	return &Server{sdk: sdk}
}

// listen listens for connections.
func (s *Server) Serve(address string) error {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
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
	return s.instance.Serve(ln)
}

// Close gracefully closes the server.
func (s *Server) Close() {
	s.instance.GracefulStop()
}

// register all server
func (s *Server) register() {
	coreapi.RegisterCoreServer(s.instance, core.NewServer(s.sdk))

	protobuf_api.RegisterEventServer(s.instance, api.NewEventServer(s.sdk))
	protobuf_api.RegisterExecutionServer(s.instance, api.NewExecutionServer(s.sdk))
	protobuf_api.RegisterInstanceServer(s.instance, api.NewInstanceServer(s.sdk))
	protobuf_api.RegisterServiceServer(s.instance, api.NewServiceServer(s.sdk))

	reflection.Register(s.instance)
}
