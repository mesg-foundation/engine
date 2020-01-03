package executionsdk

import (
	"context"

	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/api"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	runnersdk "github.com/mesg-foundation/engine/sdk/runner"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
)

// SDK is the execution sdk.
type SDK struct {
	client *cosmos.Client

	serviceSDK  *servicesdk.SDK
	instanceSDK *instancesdk.SDK
	runnerSDK   *runnersdk.SDK
}

// New returns the execution sdk.
func New(client *cosmos.Client, serviceSDK *servicesdk.SDK, instanceSDK *instancesdk.SDK, runnerSDK *runnersdk.SDK) *SDK {
	sdk := &SDK{
		client:      client,
		serviceSDK:  serviceSDK,
		instanceSDK: instanceSDK,
		runnerSDK:   runnerSDK,
	}
	return sdk
}

// Create creates a new execution.
func (s *SDK) Create(req *api.CreateExecutionRequest) (*execution.Execution, error) {
	m.Created.Add(1)
	acc, err := s.client.GetAccount()
	if err != nil {
		return nil, err
	}
	msg := newMsgCreateExecution(req, acc.GetAddress())
	tx, err := s.client.BuildAndBroadcastMsg(msg)
	m.Signed.Add(1)
	if err != nil {
		return nil, err
	}
	return s.Get(tx.Data)
}

// Update updates a execution.
func (s *SDK) Update(req *api.UpdateExecutionRequest) (*execution.Execution, error) {
	acc, err := s.client.GetAccount()
	if err != nil {
		return nil, err
	}
	msg := newMsgUpdateExecution(req, acc.GetAddress())
	tx, err := s.client.BuildAndBroadcastMsg(msg)
	if err != nil {
		return nil, err
	}
	return s.Get(tx.Data)
}

// Get returns the execution that matches given hash.
func (s *SDK) Get(hash hash.Hash) (*execution.Execution, error) {
	var execution execution.Execution
	if err := s.client.Query("custom/"+backendName+"/get/"+hash.String(), nil, &execution); err != nil {
		return nil, err
	}
	return &execution, nil
}

// List returns all executions.
func (s *SDK) List() ([]*execution.Execution, error) {
	var executions []*execution.Execution
	if err := s.client.Query("custom/"+backendName+"/list", nil, &executions); err != nil {
		return nil, err
	}
	return executions, nil
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
