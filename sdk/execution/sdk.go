package executionsdk

import (
	"context"
	"fmt"

	"github.com/mesg-foundation/engine/cosmos"
	executionpb "github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/api"
	runnersdk "github.com/mesg-foundation/engine/sdk/runner"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	"github.com/mesg-foundation/engine/x/execution"
)

// SDK is the execution sdk.
type SDK struct {
	client *cosmos.Client

	serviceSDK *servicesdk.SDK
	runnerSDK  *runnersdk.SDK
}

// New returns the execution sdk.
func New(client *cosmos.Client, serviceSDK *servicesdk.SDK, runnerSDK *runnersdk.SDK) *SDK {
	sdk := &SDK{
		client:     client,
		serviceSDK: serviceSDK,
		runnerSDK:  runnerSDK,
	}
	return sdk
}

// Create creates a new execution.
func (s *SDK) Create(req *api.CreateExecutionRequest) (*executionpb.Execution, error) {
	execution.M.Created.Add(1)
	acc, err := s.client.GetAccount()
	if err != nil {
		return nil, err
	}
	msg := execution.NewMsgCreateExecution(req, acc.GetAddress())
	tx, err := s.client.BuildAndBroadcastMsg(msg)
	if err != nil {
		return nil, err
	}
	return s.Get(tx.Data)
}

// Update updates a execution.
func (s *SDK) Update(req *api.UpdateExecutionRequest) (*executionpb.Execution, error) {
	execution.M.Updated.Add(1)
	acc, err := s.client.GetAccount()
	if err != nil {
		return nil, err
	}
	msg := execution.NewMsgUpdateExecution(req, acc.GetAddress())
	tx, err := s.client.BuildAndBroadcastMsg(msg)
	if err != nil {
		return nil, err
	}
	return s.Get(tx.Data)
}

// Get returns the execution that matches given hash.
func (s *SDK) Get(hash hash.Hash) (*executionpb.Execution, error) {
	var e executionpb.Execution
	if err := s.client.QueryJSON(fmt.Sprintf("custom/%s/%s/%s", execution.QuerierRoute, execution.QueryGetExecution, hash.String()), nil, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

// List returns all executions.
func (s *SDK) List() ([]*executionpb.Execution, error) {
	var es []*executionpb.Execution
	if err := s.client.QueryJSON(fmt.Sprintf("custom/%s/%s", execution.QuerierRoute, execution.QueryListExecution), nil, &es); err != nil {
		return nil, err
	}
	return es, nil
}

// Stream returns execution that matches given hash.
func (s *SDK) Stream(ctx context.Context, req *api.StreamExecutionRequest) (chan *executionpb.Execution, chan error, error) {
	if err := req.Filter.Validate(); err != nil {
		return nil, nil, err
	}

	stream, serrC, err := s.client.Stream(ctx, cosmos.EventModuleQuery(execution.ModuleName))
	if err != nil {
		return nil, nil, err
	}

	execC := make(chan *executionpb.Execution)
	errC := make(chan error)
	go func() {
	loop:
		for {
			select {
			case hash := <-stream:
				exec, err := s.Get(hash)
				if err != nil {
					errC <- err
					break
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
