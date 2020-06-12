package grpc

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_kit "github.com/grpc-ecosystem/go-grpc-middleware/logging/kit"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/mesg-foundation/engine/event/publisher"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/ext/xvalidator"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
	types "github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/runner"
	"github.com/mesg-foundation/engine/server/grpc/orchestrator"
	grpcrunner "github.com/mesg-foundation/engine/server/grpc/runner"
	"github.com/mesg-foundation/engine/service"
	tmlog "github.com/tendermint/tendermint/libs/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// Store is the interface to implement to fetch data.
type Store interface {
	// FetchService returns a service from its hash.
	FetchService(ctx context.Context, hash hash.Hash) (*service.Service, error)

	// FetchInstance returns an instance from its hash.
	FetchInstance(ctx context.Context, hash hash.Hash) (*instance.Instance, error)

	// CreateExecution creates an execution.
	CreateExecution(ctx context.Context, taskKey string, inputs *types.Struct, tags []string, parentHash hash.Hash, eventHash hash.Hash, processHash hash.Hash, nodeKey string, executorHash hash.Hash) (hash.Hash, error)

	// SubscribeToExecutions returns a chan that will contain executions that have been created, updated, or anything.
	SubscribeToExecutions(ctx context.Context) (<-chan *execution.Execution, error)

	// SubscribeToExecutionsForRunner returns a chan that will contain executions that a specific runner must execute.
	SubscribeToExecutionsForRunner(ctx context.Context, runnerHash hash.Hash) (<-chan *execution.Execution, error)

	// FetchExecution returns one execution from its hash.
	FetchExecution(ctx context.Context, hash hash.Hash) (*execution.Execution, error)

	// UpdateExecution update an execution.
	UpdateExecution(ctx context.Context, execHash hash.Hash, start int64, stop int64, outputs *types.Struct, err string) error

	// FetchRunner returns a runner from its hash.
	FetchRunner(ctx context.Context, hash hash.Hash) (*runner.Runner, error)

	// RegisterRunner registers a new or existing runner.
	RegisterRunner(ctx context.Context, serviceHash hash.Hash, envHash hash.Hash) (hash.Hash, error)

	// DeleteRunner deletes an existing runner.
	DeleteRunner(ctx context.Context, runnerHash hash.Hash) error
}

// Server contains the server config.
type Server struct {
	instance          *grpc.Server
	store             Store
	ep                *publisher.EventPublisher
	logger            tmlog.Logger
	authorizedPubKeys []string
	cdc               *codec.Codec
}

// New returns a new gRPC server.
func New(store Store, ep *publisher.EventPublisher, logger tmlog.Logger, cdc *codec.Codec, authorizedPubKeys []string) *Server {
	return &Server{
		store:             store,
		ep:                ep,
		logger:            logger.With("module", "grpc"),
		authorizedPubKeys: authorizedPubKeys,
		cdc:               cdc,
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
			grpc_kit.StreamServerInterceptor(newTmLogger(s.logger, "Served gRPC stream")),
			grpc_prometheus.StreamServerInterceptor,
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_kit.UnaryServerInterceptor(newTmLogger(s.logger, "Served gRPC response")),
			grpc_prometheus.UnaryServerInterceptor,
			validateInterceptor,
		)),
	)
	if err := s.register(); err != nil {
		return err
	}
	grpc_prometheus.Register(s.instance)
	s.logger.Info("Server listens on " + ln.Addr().String())
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

	grpcrunner.RegisterRunnerServer(s.instance, grpcrunner.NewServer(s.store, s.ep, tokenToRunnerHash, s.logger))

	authorizer, err := orchestrator.NewAuthorizer(s.cdc, s.authorizedPubKeys)
	if err != nil {
		return err
	}
	orchestrator.RegisterEventServer(s.instance, orchestrator.NewEventServer(s.ep, authorizer))
	orchestrator.RegisterExecutionServer(s.instance, orchestrator.NewExecutionServer(s.store, authorizer))
	orchestrator.RegisterRunnerServer(s.instance, orchestrator.NewRunnerServer(s.store, tokenToRunnerHash, authorizer))

	reflection.Register(s.instance)

	return nil
}
