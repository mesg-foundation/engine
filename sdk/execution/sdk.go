package executionsdk

import (
	"fmt"

	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/api"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	runnersdk "github.com/mesg-foundation/engine/sdk/runner"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	"github.com/tendermint/tendermint/mempool"
)

// SDK is the execution sdk.
type SDK struct {
	client *cosmos.Client
	kb     *cosmos.Keybase

	serviceSDK  *servicesdk.SDK
	instanceSDK *instancesdk.SDK
	runnerSDK   *runnersdk.SDK
}

// New returns the execution sdk.
func New(client *cosmos.Client, kb *cosmos.Keybase, serviceSDK *servicesdk.SDK, instanceSDK *instancesdk.SDK, runnerSDK *runnersdk.SDK) *SDK {
	sdk := &SDK{
		client:      client,
		kb:          kb,
		serviceSDK:  serviceSDK,
		instanceSDK: instanceSDK,
		runnerSDK:   runnerSDK,
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
		if err == mempool.ErrTxInCache {
			return nil, fmt.Errorf("execution already exists: %w", err)
		}
		return nil, err
	}
	return s.Get(tx.Data)
}

// Update updates a execution.
func (s *SDK) Update(req *api.UpdateExecutionRequest, accountName, accountPassword string) (*execution.Execution, error) {
	acc, err := s.kb.Get(accountName)
	if err != nil {
		return nil, err
	}
	msg := newMsgUpdateExecution(req, acc.GetAddress())
	tx, err := s.client.BuildAndBroadcastMsg(msg, accountName, accountPassword)
	if err != nil {
		if err == mempool.ErrTxInCache {
			return nil, fmt.Errorf("execution already exists: %w", err)
		}
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
func (s *SDK) Stream(req *api.StreamExecutionRequest) (chan *execution.Execution, error) {
	if err := req.Filter.Validate(); err != nil {
		return nil, err
	}

	stream, err := s.client.Stream(cosmos.EventModuleQuery(backendName))
	if err != nil {
		return nil, err
	}

	execC := make(chan *execution.Execution)
	go func() {
		for hash := range stream {
			exec, err := s.Get(hash)
			if err != nil {
				// or panic(err) - grpc api do not support
				// return the errors on the stream for now
				// so besieds logging the error, it not
				// much we can do here.
				continue
			}
			if req.Filter.Match(exec) {
				execC <- exec
			}
		}
		close(execC)
	}()
	return execC, nil
}
