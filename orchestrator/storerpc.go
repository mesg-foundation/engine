package orchestrator

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/ext/xstrings"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/runner"
	executionmodule "github.com/mesg-foundation/engine/x/execution"
	processmodule "github.com/mesg-foundation/engine/x/process"
	runnermodule "github.com/mesg-foundation/engine/x/runner"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

type StoreRPC struct {
	rpc    *cosmos.RPC
	logger tmlog.Logger
}

func NewStoreRPC(rpc *cosmos.RPC, logger tmlog.Logger) *StoreRPC {
	return &StoreRPC{
		rpc:    rpc,
		logger: logger,
	}
}

func (s *StoreRPC) FetchProcesses(ctx context.Context) ([]*process.Process, error) {
	var processes []*process.Process
	route := fmt.Sprintf("custom/%s/%s", processmodule.QuerierRoute, processmodule.QueryList)
	if err := s.rpc.QueryJSON(route, nil, &processes); err != nil {
		return nil, err
	}
	return processes, nil
}

func (s *StoreRPC) FetchExecution(ctx context.Context, hash hash.Hash) (*execution.Execution, error) {
	var exec *execution.Execution
	route := fmt.Sprintf("custom/%s/%s/%s", executionmodule.QuerierRoute, executionmodule.QueryGet, hash)
	if err := s.rpc.QueryJSON(route, nil, &exec); err != nil {
		return nil, err
	}
	return exec, nil
}

func (s *StoreRPC) FetchRunners(ctx context.Context, instanceHash hash.Hash) ([]*runner.Runner, error) {
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

func (s *StoreRPC) CreateExecution(ctx context.Context, taskKey string, inputs *types.Struct, tags []string, parentHash hash.Hash, eventHash hash.Hash, processHash hash.Hash, nodeKey string, executorHash hash.Hash) (hash.Hash, error) {
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

func (s *StoreRPC) SubscribeToNewCompletedExecutions(ctx context.Context) (<-chan *execution.Execution, error) {
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
