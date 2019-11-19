package executionsdk

import (
	"fmt"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/runner"
	accountsdk "github.com/mesg-foundation/engine/sdk/account"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	runnersdk "github.com/mesg-foundation/engine/sdk/runner"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	"github.com/mesg-foundation/engine/x/xstrings"
	"github.com/tendermint/tendermint/mempool"
)

// SDK is the execution sdk.
type SDK struct {
	client      *cosmos.Client
	accountSDK  *accountsdk.SDK
	serviceSDK  *servicesdk.SDK
	instanceSDK *instancesdk.SDK
	runnerSDK   *runnersdk.SDK
}

// New returns the execution sdk.
func New(client *cosmos.Client, accountSDK *accountsdk.SDK, serviceSDK *servicesdk.SDK, instanceSDK *instancesdk.SDK, runnerSDK *runnersdk.SDK) *SDK {
	sdk := &SDK{
		client:      client,
		accountSDK:  accountSDK,
		serviceSDK:  serviceSDK,
		instanceSDK: instanceSDK,
		runnerSDK:   runnerSDK,
	}
	return sdk
}

// Create creates a new execution.
func (s *SDK) Create(req *api.CreateExecutionRequest, accountName, accountPassword string) (*execution.Execution, error) {
	account, err := s.accountSDK.Get(accountName)
	if err != nil {
		return nil, err
	}
	signer, err := cosmostypes.AccAddressFromBech32(account.Address)
	if err != nil {
		return nil, err
	}
	msg := newMsgCreateExecution(req, signer)
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
	account, err := s.accountSDK.Get(accountName)
	if err != nil {
		return nil, err
	}
	executor, err := cosmostypes.AccAddressFromBech32(account.Address)
	if err != nil {
		return nil, err
	}
	msg := newMsgUpdateExecution(req, executor)
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
func (s *SDK) Stream(req *api.StreamExecutionRequest) (<-chan *execution.Execution, func() error, error) {
	if err := s.validateFilter(req.Filter); err != nil {
		return nil, func() error { return nil }, err
	}
	stream, closer, err := s.client.Stream(cosmos.EventModuleQuery(backendName))
	if err != nil {
		return nil, closer, err
	}
	execChan := make(chan *execution.Execution)
	go func() {
		defer close(execChan)
		for hash := range stream {
			exec, err := s.Get(hash)
			if err != nil {
				// TODO: remove panic
				panic(err)
			}
			if match(req.Filter, exec) {
				execChan <- exec
			}
		}
	}()
	return execChan, closer, nil
}

func (s *SDK) validateFilter(f *api.StreamExecutionRequest_Filter) error {
	if f == nil {
		return nil
	}
	var err error
	var run *runner.Runner
	if !f.ExecutorHash.IsZero() {
		if run, err = s.runnerSDK.Get(f.ExecutorHash); err != nil {
			return err
		}
	}
	var inst *instance.Instance
	if !f.InstanceHash.IsZero() {
		if inst, err = s.instanceSDK.Get(f.InstanceHash); err != nil {
			return err
		}
	}
	if (f.TaskKey == "" || f.TaskKey == "*") || (inst == nil && run == nil) {
		return nil
	}
	// check task key if at least instance or runner is set
	if inst == nil && run != nil {
		inst, err = s.instanceSDK.Get(run.InstanceHash)
		if err != nil {
			return err
		}
	}
	srv, err := s.serviceSDK.Get(inst.ServiceHash)
	if err != nil {
		return err
	}
	if _, err := srv.GetTask(f.TaskKey); err != nil {
		return err
	}
	return nil
}

// Match matches an execution against a filter.
func match(f *api.StreamExecutionRequest_Filter, e *execution.Execution) bool {
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
