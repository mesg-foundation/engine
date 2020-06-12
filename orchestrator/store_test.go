package orchestrator

import (
	"context"
	"errors"
	"fmt"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cskr/pubsub"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/types"
	"github.com/mesg-foundation/engine/runner"
	"github.com/mesg-foundation/engine/service"
)

const (
	pubsubExecTopic = "exec"
)

type storeTest struct {
	processes  []*process.Process
	executions []*execution.Execution
	services   []*service.Service
	instances  []*instance.Instance
	runners    []*runner.Runner

	pubsub             *pubsub.PubSub
	runnerOwner        string
	procPaymentAddress sdktypes.AccAddress
}

func init() {
	cosmos.InitConfig()
}

func newStoreTest() (*storeTest, error) {
	procPaymentAddress, err := sdktypes.AccAddressFromBech32("mesg1t9h20sn3lk2jdnak5eea4lkqxkpwyfaadtqk4t")
	if err != nil {
		return nil, err
	}
	return &storeTest{
		processes:  make([]*process.Process, 0),
		executions: make([]*execution.Execution, 0),
		services:   make([]*service.Service, 0),
		instances:  make([]*instance.Instance, 0),
		runners:    make([]*runner.Runner, 0),

		pubsub:             pubsub.New(0),
		runnerOwner:        "mesg1s6mqusxaq93d70jeekqehg7aepwt7zs306ctq7",
		procPaymentAddress: procPaymentAddress,
	}, nil
}

func (s *storeTest) CreateProcess(ctx context.Context, name string, nodes []*process.Process_Node, edges []*process.Process_Edge) (hash.Hash, error) {
	proc, err := process.New(name, nodes, edges, s.procPaymentAddress)
	if err != nil {
		return nil, err
	}
	s.processes = append(s.processes, proc)
	return proc.Hash, nil
}

func (s *storeTest) CreateService(ctx context.Context, sid, name, description string, configuration service.Service_Configuration, tasks []*service.Service_Task, events []*service.Service_Event, dependencies []*service.Service_Dependency, repository, source string) (hash.Hash, error) {
	srv, err := service.New(sid, name, description, configuration, tasks, events, dependencies, repository, source)
	if err != nil {
		return nil, err
	}
	s.services = append(s.services, srv)
	return srv.Hash, nil
}

func (s *storeTest) FetchProcesses(ctx context.Context) ([]*process.Process, error) {
	return s.processes, nil
}

// FetchExecution returns an execution from its hash.
func (s *storeTest) FetchExecution(ctx context.Context, hash hash.Hash) (*execution.Execution, error) {
	for _, exec := range s.executions {
		if exec.Hash.Equal(hash) {
			return exec, nil
		}
	}
	return nil, fmt.Errorf("execution %q not found", hash)
}

// FetchService returns a service from its hash.
func (s *storeTest) FetchService(ctx context.Context, hash hash.Hash) (*service.Service, error) {
	for _, srv := range s.services {
		if srv.Hash.Equal(hash) {
			return srv, nil
		}
	}
	return nil, fmt.Errorf("service %q not found", hash)
}

// FetchInstance returns an instance from its hash.
func (s *storeTest) FetchInstance(ctx context.Context, hash hash.Hash) (*instance.Instance, error) {
	for _, inst := range s.instances {
		if inst.Hash.Equal(hash) {
			return inst, nil
		}
	}
	return nil, fmt.Errorf("instance %q not found", hash)
}

// FetchRunner returns a runner from its hash.
func (s *storeTest) FetchRunner(ctx context.Context, hash hash.Hash) (*runner.Runner, error) {
	for _, run := range s.runners {
		if run.Hash.Equal(hash) {
			return run, nil
		}
	}
	return nil, fmt.Errorf("runner %q not found", hash)
}

// FetchRunners returns all runners of an instance.
func (s *storeTest) FetchRunners(ctx context.Context, instanceHash hash.Hash) ([]*runner.Runner, error) {
	executors := make([]*runner.Runner, 0)
	for _, run := range s.runners {
		if run.InstanceHash.Equal(instanceHash) {
			executors = append(executors, run)
		}
	}
	return executors, nil
}

// CreateExecution creates an execution.
func (s *storeTest) CreateExecution(ctx context.Context, taskKey string, inputs *types.Struct, tags []string, parentHash hash.Hash, eventHash hash.Hash, processHash hash.Hash, nodeKey string, executorHash hash.Hash) (hash.Hash, error) {
	run, err := s.FetchRunner(ctx, executorHash)
	if err != nil {
		return nil, err
	}
	exec, err := execution.New(
		processHash,
		run.InstanceHash,
		parentHash,
		eventHash,
		nodeKey,
		taskKey,
		inputs,
		tags,
		executorHash,
	)
	if err != nil {
		return nil, err
	}
	if execExist, _ := s.FetchExecution(ctx, exec.Hash); execExist != nil {
		return nil, fmt.Errorf("execution %q already exists", exec.Hash)
	}
	if err := exec.Execute(); err != nil {
		return nil, err
	}
	s.executions = append(s.executions, exec)
	s.pubsub.Pub(exec, pubsubExecTopic)
	return exec.Hash, nil
}

// UpdateExecution update an execution.
func (s *storeTest) UpdateExecution(ctx context.Context, execHash hash.Hash, start int64, stop int64, outputs *types.Struct, outputErr string) error {
	exec, err := s.FetchExecution(ctx, execHash)
	if err != nil {
		return err
	}
	exec.Start = start
	exec.Stop = stop
	exec.Price = "10"
	if outputs != nil {
		if err := exec.Complete(outputs); err != nil {
			return err
		}
	} else {
		if err := exec.Fail(errors.New(outputErr)); err != nil {
			return err
		}
	}
	s.pubsub.Pub(exec, pubsubExecTopic)
	return nil
}

// SubscribeToNewCompletedExecutions returns a chan that will contain newly completed execution.
func (s *storeTest) SubscribeToNewCompletedExecutions(ctx context.Context) (<-chan *execution.Execution, error) {
	execChan := make(chan *execution.Execution)
	c := s.pubsub.Sub(pubsubExecTopic)
	go func() {
		for {
			select {
			case v := <-c:
				if exec, ok := v.(*execution.Execution); ok {
					if exec.Status == execution.Status_Completed {
						execChan <- exec
					}
				}
			case <-ctx.Done():
				s.pubsub.Unsub(c, pubsubExecTopic)
				close(execChan)
				return
			}
		}
	}()
	return execChan, nil
}

// SubscribeToExecutions returns a chan that will contain executions that have been created, updated, or anything.
func (s *storeTest) SubscribeToExecutions(ctx context.Context) (<-chan *execution.Execution, error) {
	execChan := make(chan *execution.Execution)
	c := s.pubsub.Sub(pubsubExecTopic)
	go func() {
		for {
			select {
			case v := <-c:
				if exec, ok := v.(*execution.Execution); ok {
					execChan <- exec
				}
			case <-ctx.Done():
				s.pubsub.Unsub(c, pubsubExecTopic)
				close(execChan)
				return
			}
		}
	}()
	return execChan, nil
}

// SubscribeToExecutionsForRunner returns a chan that will contain executions that a specific runner must execute.
func (s *storeTest) SubscribeToExecutionsForRunner(ctx context.Context, runnerHash hash.Hash) (<-chan *execution.Execution, error) {
	execChan := make(chan *execution.Execution)
	c := s.pubsub.Sub(pubsubExecTopic)
	go func() {
		for {
			select {
			case v := <-c:
				if exec, ok := v.(*execution.Execution); ok {
					if exec.Status == execution.Status_InProgress && exec.ExecutorHash.Equal(runnerHash) {
						execChan <- exec
					}
				}
			case <-ctx.Done():
				s.pubsub.Unsub(c, pubsubExecTopic)
				close(execChan)
				return
			}
		}
	}()
	return execChan, nil

}

// RegisterRunner registers a new or existing runner.
func (s *storeTest) RegisterRunner(ctx context.Context, serviceHash hash.Hash, envHash hash.Hash) (hash.Hash, error) {
	inst, err := instance.New(serviceHash, envHash)
	if err != nil {
		return nil, err
	}
	if instExist, _ := s.FetchInstance(ctx, inst.Hash); instExist == nil {
		s.instances = append(s.instances, inst)
	}

	run, err := runner.New(s.runnerOwner, inst.Hash)
	if err != nil {
		return nil, err
	}
	if runExist, _ := s.FetchRunner(ctx, run.Hash); runExist != nil {
		return nil, fmt.Errorf("runner %q already exists", run.Hash)
	}
	s.runners = append(s.runners, run)
	return run.Hash, nil
}

// DeleteRunner deletes an existing runner.
func (s *storeTest) DeleteRunner(ctx context.Context, runnerHash hash.Hash) error {
	index := -1
	for i, run := range s.runners {
		if run.Hash.Equal(runnerHash) {
			index = i
			break
		}
	}
	if index == -1 {
		return fmt.Errorf("runner %q not found", runnerHash)
	}
	s.runners = append(s.runners[:index], s.runners[index+1:]...)
	return nil
}
