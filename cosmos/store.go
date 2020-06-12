package cosmos

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/ext/xstrings"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/runner"
	"github.com/mesg-foundation/engine/service"
	executionmodule "github.com/mesg-foundation/engine/x/execution"
	instancemodule "github.com/mesg-foundation/engine/x/instance"
	processmodule "github.com/mesg-foundation/engine/x/process"
	runnermodule "github.com/mesg-foundation/engine/x/runner"
	servicemodule "github.com/mesg-foundation/engine/x/service"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

// Store is an implementation of the orchestrator.Store interface using cosmos rpc.
type Store struct {
	rpc    *RPC
	logger tmlog.Logger
}

// NewStore returns a new implementation of orchestrator.Store using cosmos rpc.
func NewStore(rpc *RPC, logger tmlog.Logger) *Store {
	return &Store{
		rpc:    rpc,
		logger: logger,
	}
}

// FetchProcesses returns all processes.
func (s *Store) FetchProcesses(ctx context.Context) ([]*process.Process, error) {
	var processes []*process.Process
	route := fmt.Sprintf("custom/%s/%s", processmodule.QuerierRoute, processmodule.QueryList)
	if err := s.rpc.QueryJSON(route, nil, &processes); err != nil {
		return nil, err
	}
	return processes, nil
}

// FetchExecution returns an execution from its hash.
func (s *Store) FetchExecution(ctx context.Context, hash hash.Hash) (*execution.Execution, error) {
	var exec *execution.Execution
	route := fmt.Sprintf("custom/%s/%s/%s", executionmodule.QuerierRoute, executionmodule.QueryGet, hash)
	if err := s.rpc.QueryJSON(route, nil, &exec); err != nil {
		return nil, err
	}
	return exec, nil
}

// FetchService returns a service from its hash.
func (s *Store) FetchService(ctx context.Context, hash hash.Hash) (*service.Service, error) {
	var srv *service.Service
	route := fmt.Sprintf("custom/%s/%s/%s", servicemodule.QuerierRoute, servicemodule.QueryGet, hash)
	if err := s.rpc.QueryJSON(route, nil, &srv); err != nil {
		return nil, err
	}
	return srv, nil
}

// FetchInstance returns an instance from its hash.
func (s *Store) FetchInstance(ctx context.Context, hash hash.Hash) (*instance.Instance, error) {
	var inst *instance.Instance
	route := fmt.Sprintf("custom/%s/%s/%s", instancemodule.QuerierRoute, instancemodule.QueryGet, hash)
	if err := s.rpc.QueryJSON(route, nil, &inst); err != nil {
		return nil, err
	}
	return inst, nil
}

// FetchRunner returns a runner from its hash.
func (s *Store) FetchRunner(ctx context.Context, hash hash.Hash) (*runner.Runner, error) {
	var run *runner.Runner
	route := fmt.Sprintf("custom/%s/%s/%s", runnermodule.QuerierRoute, runnermodule.QueryGet, hash)
	if err := s.rpc.QueryJSON(route, nil, &run); err != nil {
		return nil, err
	}
	return run, nil
}

// FetchRunners returns all runners of an instance.
func (s *Store) FetchRunners(ctx context.Context, instanceHash hash.Hash) ([]*runner.Runner, error) {
	var runners []*runner.Runner
	route := fmt.Sprintf("custom/%s/%s", runnermodule.QuerierRoute, runnermodule.QueryList)
	if err := s.rpc.QueryJSON(route, nil, &runners); err != nil {
		return nil, err
	}
	executors := make([]*runner.Runner, 0)
	for _, run := range runners {
		if run.InstanceHash.Equal(instanceHash) {
			executors = append(executors, run)
		}
	}
	return executors, nil
}

// CreateExecution creates an execution.
func (s *Store) CreateExecution(ctx context.Context, taskKey string, inputs *types.Struct, tags []string, parentHash hash.Hash, eventHash hash.Hash, processHash hash.Hash, nodeKey string, executorHash hash.Hash) (hash.Hash, error) {
	acc, err := s.rpc.GetAccount()
	if err != nil {
		return nil, err
	}
	msg := executionmodule.MsgCreate{
		Signer:       acc.GetAddress(),
		ProcessHash:  processHash,
		EventHash:    eventHash,
		ParentHash:   parentHash,
		NodeKey:      nodeKey,
		TaskKey:      taskKey,
		Inputs:       inputs,
		ExecutorHash: executorHash,
		Tags:         tags,
	}
	res, err := s.rpc.BuildAndBroadcastMsg(msg)
	if err != nil {
		return nil, err
	}
	return hash.DecodeFromBytes(res.Data)
}

// UpdateExecution update an execution.
func (s *Store) UpdateExecution(ctx context.Context, execHash hash.Hash, start int64, stop int64, outputs *types.Struct, outputErr string) error {
	acc, err := s.rpc.GetAccount()
	if err != nil {
		return err
	}
	msg := executionmodule.MsgUpdate{
		Executor: acc.GetAddress(),
		Hash:     execHash,
		Start:    start,
		Stop:     stop,
	}
	if outputs != nil {
		msg.Result = &executionmodule.MsgUpdateOutputs{
			Outputs: outputs,
		}
	} else {
		msg.Result = &executionmodule.MsgUpdateError{
			Error: outputErr,
		}
	}
	if _, err := s.rpc.BuildAndBroadcastMsg(msg); err != nil {
		return err
	}
	return nil
}

// SubscribeToNewCompletedExecutions returns a chan that will contain newly completed execution.
func (s *Store) SubscribeToNewCompletedExecutions(ctx context.Context) (<-chan *execution.Execution, error) {
	subscriber := xstrings.RandASCIILetters(8)
	query := fmt.Sprintf("%s.%s EXISTS AND %s.%s='%s'",
		executionmodule.EventType, executionmodule.AttributeKeyHash,
		executionmodule.EventType, sdk.AttributeKeyAction, executionmodule.AttributeActionCompleted,
	)
	eventStream, err := s.rpc.Subscribe(ctx, subscriber, query, 0)
	if err != nil {
		return nil, err
	}
	executionStream := make(chan *execution.Execution)
	go func() {
	loop:
		for {
			select {
			case event := <-eventStream:
				// get the index of the action=completed attributes
				attrKeyActionCreated := fmt.Sprintf("%s.%s", executionmodule.EventType, sdk.AttributeKeyAction)
				attrIndexes := make([]int, 0)
				for index, attr := range event.Events[attrKeyActionCreated] {
					if attr == executionmodule.AttributeActionCompleted {
						attrIndexes = append(attrIndexes, index)
					}
				}
				// iterate only on the index of attribute hash where action=completed
				attrKeyHash := fmt.Sprintf("%s.%s", executionmodule.EventType, executionmodule.AttributeKeyHash)
				for _, index := range attrIndexes {
					attr := event.Events[attrKeyHash][index]
					hash, err := hash.Decode(attr)
					if err != nil {
						s.logger.Error(err.Error())
						continue
					}
					var exec *execution.Execution
					route := fmt.Sprintf("custom/%s/%s/%s", executionmodule.QuerierRoute, executionmodule.QueryGet, hash)
					if err := s.rpc.QueryJSON(route, nil, &exec); err != nil {
						s.logger.Error(err.Error())
						continue
					}
					executionStream <- exec
				}
			case <-ctx.Done():
				break loop
			}
		}
		if err := s.rpc.Unsubscribe(context.Background(), subscriber, query); err != nil {
			s.logger.Error(err.Error())
		}
	}()
	return executionStream, nil
}

// SubscribeToExecutions returns a chan that will contain executions that have been created, updated, or anything.
func (s *Store) SubscribeToExecutions(ctx context.Context) (<-chan *execution.Execution, error) {
	subscriber := xstrings.RandASCIILetters(8)
	query := fmt.Sprintf("%s.%s EXISTS", executionmodule.EventType, executionmodule.AttributeKeyHash)
	eventStream, err := s.rpc.Subscribe(ctx, subscriber, query, 0)
	if err != nil {
		return nil, err
	}
	executionStream := make(chan *execution.Execution)
	// listen to event stream
	go func() {
	loop:
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
						s.logger.Error(err.Error())
						continue
					}
					var exec *execution.Execution
					route := fmt.Sprintf("custom/%s/%s/%s", executionmodule.QuerierRoute, executionmodule.QueryGet, hash)
					if err := s.rpc.QueryJSON(route, nil, &exec); err != nil {
						s.logger.Error(err.Error())
						continue
					}
					executionStream <- exec
				}
			case <-ctx.Done():
				break loop
			}
		}
		if err := s.rpc.Unsubscribe(context.Background(), subscriber, query); err != nil {
			s.logger.Error(err.Error())
		}
	}()
	return executionStream, nil
}

// SubscribeToExecutionsForRunner returns a chan that will contain executions that a specific runner must execute.
func (s *Store) SubscribeToExecutionsForRunner(ctx context.Context, runnerHash hash.Hash) (<-chan *execution.Execution, error) {
	subscriber := xstrings.RandASCIILetters(8)
	query := fmt.Sprintf("%s.%s='%s' AND %s.%s='%s'",
		executionmodule.EventType, executionmodule.AttributeKeyExecutor, runnerHash.String(),
		executionmodule.EventType, sdk.AttributeKeyAction, executionmodule.AttributeActionCreated,
	)
	eventStream, err := s.rpc.Subscribe(ctx, subscriber, query, 0)
	if err != nil {
		return nil, err
	}
	executionStream := make(chan *execution.Execution)
	// listen to event stream
	go func() {
	loop:
		for {
			select {
			case event := <-eventStream:
				// get the index of the action=created attributes
				attrKeyActionCreated := fmt.Sprintf("%s.%s", executionmodule.EventType, sdk.AttributeKeyAction)
				attrIndexes := make([]int, 0)
				for index, attr := range event.Events[attrKeyActionCreated] {
					if attr == executionmodule.AttributeActionCreated {
						attrIndexes = append(attrIndexes, index)
					}
				}
				// iterate only on the index of attribute hash where action=created
				attrKeyHash := fmt.Sprintf("%s.%s", executionmodule.EventType, executionmodule.AttributeKeyHash)
				for _, index := range attrIndexes {
					attr := event.Events[attrKeyHash][index]
					hash, err := hash.Decode(attr)
					if err != nil {
						s.logger.Error(err.Error())
						continue
					}
					var exec *execution.Execution
					route := fmt.Sprintf("custom/%s/%s/%s", executionmodule.QuerierRoute, executionmodule.QueryGet, hash)
					if err := s.rpc.QueryJSON(route, nil, &exec); err != nil {
						s.logger.Error(err.Error())
						continue
					}
					executionStream <- exec
				}
			case <-ctx.Done():
				break loop
			}
		}
		if err := s.rpc.Unsubscribe(context.Background(), subscriber, query); err != nil {
			s.logger.Error(err.Error())
		}
	}()
	return executionStream, nil
}

// RegisterRunner registers a new or existing runner.
func (s *Store) RegisterRunner(ctx context.Context, serviceHash hash.Hash, envHash hash.Hash) (hash.Hash, error) {
	// get engine account
	acc, err := s.rpc.GetAccount()
	if err != nil {
		return nil, err
	}

	// calculate runner hash
	inst, err := instance.New(serviceHash, envHash)
	if err != nil {
		return nil, err
	}
	run, err := runner.New(acc.GetAddress().String(), inst.Hash)
	if err != nil {
		return nil, err
	}
	runnerHash := run.Hash

	// check that runner doesn't already exist
	var runnerExist bool
	route := fmt.Sprintf("custom/%s/%s/%s", runnermodule.QuerierRoute, runnermodule.QueryExist, runnerHash)
	if err := s.rpc.QueryJSON(route, nil, &runnerExist); err != nil {
		return nil, err
	}

	// only broadcast if runner doesn't exist
	if !runnerExist {
		tx, err := s.rpc.BuildAndBroadcastMsg(runnermodule.MsgCreate{
			Owner:       acc.GetAddress(),
			ServiceHash: serviceHash,
			EnvHash:     envHash,
		})
		if err != nil {
			return nil, err
		}
		runnerHashCreated, err := hash.DecodeFromBytes(tx.Data)
		if err != nil {
			return nil, err
		}
		if !runnerHashCreated.Equal(runnerHash) {
			// delete wrong runner
			_, err := s.rpc.BuildAndBroadcastMsg(runnermodule.MsgDelete{
				Owner: acc.GetAddress(),
				Hash:  runnerHashCreated,
			})
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("runner hash created is not expected: got %q, expect %q", runnerHashCreated, runnerHash)
		}
	}

	return runnerHash, nil
}

// DeleteRunner deletes an existing runner.
func (s *Store) DeleteRunner(ctx context.Context, runnerHash hash.Hash) error {
	acc, err := s.rpc.GetAccount()
	if err != nil {
		return err
	}
	msg := runnermodule.MsgDelete{
		Owner: acc.GetAddress(),
		Hash:  runnerHash,
	}
	if _, err := s.rpc.BuildAndBroadcastMsg(msg); err != nil {
		return err
	}
	return nil
}
