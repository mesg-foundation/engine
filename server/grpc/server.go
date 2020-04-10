package grpc

import (
	"context"
	"net"
	"sync"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/event/publisher"
	"github.com/mesg-foundation/engine/ext/xvalidator"
	"github.com/mesg-foundation/engine/server/grpc/orchestrator"
	"github.com/mesg-foundation/engine/server/grpc/runner"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// Server contains the server config.
type Server struct {
	instance          *grpc.Server
	rpc               *cosmos.RPC
	ep                *publisher.EventPublisher
	authorizedPubKeys []string
	accName           string
	accPassword       string
}

// New returns a new gRPC server.
func New(rpc *cosmos.RPC, ep *publisher.EventPublisher, authorizedPubKeys []string, accName, accPassword string) *Server {
	return &Server{
		rpc:               rpc,
		ep:                ep,
		authorizedPubKeys: authorizedPubKeys,
		accName:           accName,
		accPassword:       accPassword,
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
	if err := s.register(); err != nil {
		return err
	}
	grpc_prometheus.Register(s.instance)
	logrus.WithField("module", "grpc").Info("server listens on ", ln.Addr())
	return s.instance.Serve(ln)
}

// Close gracefully closes the server.
func (s *Server) Close() {
	s.instance.GracefulStop()
}

// TODO: could use github.com/grpc-ecosystem/go-grpc-middleware@v1.2.0/validator/validator.go for validating any request with a `Validate() error` function.
func validateInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := xvalidator.Struct(req); err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

// register all server
func (s *Server) register() error {
	tokenToRunnerHash := &sync.Map{}

	runner.RegisterRunnerServer(s.instance, runner.NewServer(s.rpc, s.ep, tokenToRunnerHash, s.accName, s.accPassword))

	authorizer, err := orchestrator.NewAuthorizer(s.rpc.Codec(), s.authorizedPubKeys)
	if err != nil {
		return err
	}
	orchestrator.RegisterEventServer(s.instance, orchestrator.NewEventServer(s.ep, authorizer))
	orchestrator.RegisterExecutionServer(s.instance, orchestrator.NewExecutionServer(s.rpc, authorizer, s.accName, s.accPassword))
	orchestrator.RegisterRunnerServer(s.instance, orchestrator.NewRunnerServer(s.rpc, tokenToRunnerHash, authorizer, s.accName, s.accPassword))

	reflection.Register(s.instance)

	return nil
}
