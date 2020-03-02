package grpc

import (
	"context"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/event/publisher"
	"github.com/mesg-foundation/engine/ext/xvalidator"
	protobuf_api "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/runner/builder"
	"github.com/mesg-foundation/engine/server/grpc/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// Server contains the server config.
type Server struct {
	instance  *grpc.Server
	mc        *cosmos.ModuleClient
	ep        *publisher.EventPublisher
	b         *builder.Builder
	execPrice string
}

// New returns a new gRPC server.
func New(mc *cosmos.ModuleClient, ep *publisher.EventPublisher, b *builder.Builder, execPrice string) *Server {
	return &Server{
		mc:        mc,
		ep:        ep,
		b:         b,
		execPrice: execPrice,
	}
}

// Serve listens for connections.
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
			grpc_logrus.StreamServerInterceptor(logrus.StandardLogger().WithField("module", "grpc")),
			grpc_prometheus.StreamServerInterceptor,
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(logrus.StandardLogger().WithField("module", "grpc")),
			grpc_prometheus.UnaryServerInterceptor,
			validateInterceptor,
		)),
	)
	s.register()
	grpc_prometheus.Register(s.instance)
	logrus.WithField("module", "grpc").Info("server listens on ", ln.Addr())
	return s.instance.Serve(ln)
}

// Close gracefully closes the server.
func (s *Server) Close() {
	s.instance.GracefulStop()
}

func validateInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := xvalidator.Validate.Struct(req); err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

// register all server
func (s *Server) register() {
	protobuf_api.RegisterEventServer(s.instance, api.NewEventServer(s.ep))
	protobuf_api.RegisterExecutionServer(s.instance, api.NewExecutionServer(s.mc, s.execPrice))
	protobuf_api.RegisterInstanceServer(s.instance, api.NewInstanceServer(s.mc))
	protobuf_api.RegisterServiceServer(s.instance, api.NewServiceServer(s.mc))
	protobuf_api.RegisterProcessServer(s.instance, api.NewProcessServer(s.mc))
	protobuf_api.RegisterOwnershipServer(s.instance, api.NewOwnershipServer(s.mc))
	protobuf_api.RegisterRunnerServer(s.instance, api.NewRunnerServer(s.mc, s.b))

	reflection.Register(s.instance)
}
