package executionsdk

import (
	"context"

	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/api"
)

// SDK is the execution sdk.
type SDK struct {
	client *cosmos.Client
	kb     *cosmos.Keybase
}

// New returns the execution sdk.
func New(client *cosmos.Client, kb *cosmos.Keybase) *SDK {
	sdk := &SDK{
		client: client,
		kb:     kb,
	}
	return sdk
}

// Create creates a new execution.
func (s *SDK) Create(req *api.CreateExecutionRequest, accountName, accountPassword string) (*execution.Execution, error) {
	acc, err := s.kb.Get(accountName)
	if err != nil {
		return nil, err
	}

	msg := newMsgCreateExecution(req, acc.GetAddress())
	tx, err := s.client.BuildAndBroadcastMsg(msg, accountName, accountPassword)
	if err != nil {
		return nil, err
	}
	return s.Get(tx.Data)
}

// Get returns the execution that matches given hash.
func (s *SDK) Get(hash hash.Hash) (*execution.Execution, error) {
	var exec *execution.Execution
	if err := s.client.Query("custom/"+backendName+"/get/"+hash.String(), nil, &exec); err != nil {
		return nil, err
	}
	return exec, nil
}

// List returns all executions.
func (s *SDK) List() ([]*execution.Execution, error) {
	var execs []*execution.Execution
	if err := s.client.Query("custom/"+backendName+"/list", nil, &execs); err != nil {
		return nil, err
	}
	return execs, nil
}

// Stream returns execution that matches given hash.
func (s *SDK) Stream(ctx context.Context, req *api.StreamExecutionRequest) (chan *execution.Execution, chan error, error) {
	if err := req.Filter.Validate(); err != nil {
		return nil, nil, err
	}
	stream, serrC, err := s.client.Stream(ctx, cosmos.EventModuleQuery(backendName))
	if err != nil {
		return nil, nil, err
	}
	execC := make(chan *execution.Execution)
	errC := make(chan error)
	go func() {
	loop:
		for {
			select {
			case hash := <-stream:
				exec, err := s.Get(hash)
				if err != nil {
					errC <- err
				}
				if req.Filter.Match(exec) {
					execC <- exec
				}
			case err := <-serrC:
				errC <- err
			case <-ctx.Done():
				break loop
			}
		}
		close(errC)
		close(execC)
	}()
	return execC, errC, nil
}
