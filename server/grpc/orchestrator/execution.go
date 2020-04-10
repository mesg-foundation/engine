package orchestrator

import (
	"context"
	fmt "fmt"

	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/ext/xstrings"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/acknowledgement"
	executionmodule "github.com/mesg-foundation/engine/x/execution"
)

type executionServer struct {
	rpc  *cosmos.RPC
	auth *Authorizer
}

// NewExecutionServer creates a new Execution Server.
func NewExecutionServer(rpc *cosmos.RPC, auth *Authorizer) ExecutionServer {
	return &executionServer{
		rpc:  rpc,
		auth: auth,
	}
}

// Create creates an execution.
func (s *executionServer) Create(ctx context.Context, req *ExecutionCreateRequest) (*ExecutionCreateResponse, error) {
	// check authorization
	if err := s.auth.IsAuthorized(ctx, req); err != nil {
		return nil, err
	}

	// create execution
	acc, err := s.rpc.GetAccount()
	if err != nil {
		return nil, err
	}
	eventHash, err := hash.Random()
	if err != nil {
		return nil, err
	}
	msg := executionmodule.MsgCreate{
		Signer:       acc.GetAddress(),
		EventHash:    eventHash,
		ExecutorHash: req.ExecutorHash,
		Inputs:       req.Inputs,
		Price:        req.Price,
		Tags:         req.Tags,
		TaskKey:      req.TaskKey,
	}
	tx, err := s.rpc.BuildAndBroadcastMsg(msg)
	if err != nil {
		return nil, err
	}
	return &ExecutionCreateResponse{
		Hash: tx.Data,
	}, nil
}

// Stream returns stream of executions.
func (s *executionServer) Stream(req *ExecutionStreamRequest, stream Execution_StreamServer) error {
	// check authorization
	if err := s.auth.IsAuthorized(stream.Context(), req); err != nil {
		return err
	}

	// create rpc event stream
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	subscriber := xstrings.RandASCIILetters(8)
	query := fmt.Sprintf("%s.%s EXISTS", executionmodule.EventType, executionmodule.AttributeKeyHash)
	eventStream, err := s.rpc.Subscribe(ctx, subscriber, query, 0)
	defer s.rpc.Unsubscribe(context.Background(), subscriber, query)
	if err != nil {
		return err
	}
	if err := acknowledgement.SetStreamReady(stream); err != nil {
		return err
	}

	// listen to event stream
	for {
		select {
		case event := <-eventStream:
			attrHash := fmt.Sprintf("%s.%s", executionmodule.EventType, executionmodule.AttributeKeyHash)
			attrs := event.Events[attrHash]
			alreadySeeHashes := make(map[string]bool)
			for _, attr := range attrs {
				// skip already see hash. it deduplicate same execution in multiple event.
				if alreadySeeHashes[attr] {
					continue
				}
				alreadySeeHashes[attr] = true
				hash, err := hash.Decode(attr)
				if err != nil {
					return err
				}
				var exec *execution.Execution
				route := fmt.Sprintf("custom/%s/%s/%s", executionmodule.QuerierRoute, executionmodule.QueryGet, hash)
				if err := s.rpc.QueryJSON(route, nil, &exec); err != nil {
					return err
				}
				if !req.Filter.Match(exec) {
					continue
				}
				if err := stream.Send(exec); err != nil {
					return err
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}

// Match matches given execution with filter criteria.
func (f *ExecutionStreamRequest_Filter) Match(e *execution.Execution) bool {
	if f == nil {
		return true
	}
	if !f.ExecutorHash.IsZero() && !f.ExecutorHash.Equal(e.ExecutorHash) {
		return false
	}
	if !f.InstanceHash.IsZero() && !f.InstanceHash.Equal(e.InstanceHash) {
		return false
	}
	if f.TaskKey != "" && f.TaskKey != "*" && f.TaskKey != e.TaskKey {
		return false
	}

	match := len(f.Statuses) == 0
	for _, status := range f.Statuses {
		if status == e.Status {
			match = true
			break
		}
	}
	if !match {
		return false
	}
	for _, tag := range f.Tags {
		if !xstrings.SliceContains(e.Tags, tag) {
			return false
		}
	}
	return true
}
