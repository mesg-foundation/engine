package executionsdk

import (
	"errors"
	"fmt"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	processsdk "github.com/mesg-foundation/engine/sdk/process"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	"github.com/mesg-foundation/engine/x/instance"
	"github.com/mesg-foundation/engine/x/runner"
)

// Keeper holds the logic to read and write data.
type Keeper struct {
	storeKey       *cosmostypes.KVStoreKey
	serviceKeeper  *servicesdk.Keeper
	instanceKeeper instance.Keeper
	runnerKeeper   runner.Keeper
	processKeeper  *processsdk.Keeper
}

// NewKeeper initialize a new keeper.
func NewKeeper(storeKey *cosmostypes.KVStoreKey, serviceKeeper *servicesdk.Keeper, instanceKeeper instance.Keeper, runnerKeeper runner.Keeper, processKeeper *processsdk.Keeper) *Keeper {
	return &Keeper{
		storeKey:       storeKey,
		serviceKeeper:  serviceKeeper,
		instanceKeeper: instanceKeeper,
		runnerKeeper:   runnerKeeper,
		processKeeper:  processKeeper,
	}
}

// Create creates a new execution from definition.
func (k *Keeper) Create(request cosmostypes.Request, msg msgCreateExecution) (*execution.Execution, error) {
	run, err := k.runnerKeeper.Get(request, msg.Request.ExecutorHash)
	if err != nil {
		return nil, err
	}
	inst, err := k.instanceKeeper.Get(request, run.InstanceHash)
	if err != nil {
		return nil, err
	}
	srv, err := k.serviceKeeper.Get(request, inst.ServiceHash)
	if err != nil {
		return nil, err
	}
	if !msg.Request.ProcessHash.IsZero() {
		if _, err := k.processKeeper.Get(request, msg.Request.ProcessHash); err != nil {
			return nil, err
		}
	}
	if err := srv.RequireTaskInputs(msg.Request.TaskKey, msg.Request.Inputs); err != nil {
		return nil, err
	}
	exec := execution.New(
		msg.Request.ProcessHash,
		run.InstanceHash,
		msg.Request.ParentHash,
		msg.Request.EventHash,
		msg.Request.NodeKey,
		msg.Request.TaskKey,
		msg.Request.Inputs,
		msg.Request.Tags,
		msg.Request.ExecutorHash,
	)
	store := request.KVStore(k.storeKey)
	if store.Has(exec.Hash) {
		return nil, fmt.Errorf("execution %q already exists", exec.Hash)
	}
	if err := exec.Execute(); err != nil {
		return nil, err
	}
	value, err := codec.MarshalBinaryBare(exec)
	if err != nil {
		return nil, err
	}
	if !request.IsCheckTx() {
		m.InProgress.Add(1)
	}
	store.Set(exec.Hash, value)
	return exec, nil
}

// Update updates a new execution from definition.
func (k *Keeper) Update(request cosmostypes.Request, msg msgUpdateExecution) (*execution.Execution, error) {
	store := request.KVStore(k.storeKey)
	if !store.Has(msg.Request.Hash) {
		return nil, fmt.Errorf("execution %q doesn't exist", msg.Request.Hash)
	}
	var exec *execution.Execution
	if err := codec.UnmarshalBinaryBare(store.Get(msg.Request.Hash), &exec); err != nil {
		return nil, err
	}
	switch res := msg.Request.Result.(type) {
	case *api.UpdateExecutionRequest_Outputs:
		if err := k.validateExecutionOutput(request, exec.InstanceHash, exec.TaskKey, res.Outputs); err != nil {
			if err1 := exec.Failed(err); err1 != nil {
				return nil, err1
			}
		} else if err := exec.Complete(res.Outputs); err != nil {
			return nil, err
		}
	case *api.UpdateExecutionRequest_Error:
		if err := exec.Failed(errors.New(res.Error)); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("no execution result supplied")
	}
	value, err := codec.MarshalBinaryBare(exec)
	if err != nil {
		return nil, err
	}
	if !request.IsCheckTx() {
		m.Completed.Add(1)
	}
	store.Set(exec.Hash, value)
	return exec, nil
}

func (k *Keeper) validateExecutionOutput(request cosmostypes.Request, instanceHash hash.Hash, taskKey string, outputs *types.Struct) error {
	inst, err := k.instanceKeeper.Get(request, instanceHash)
	if err != nil {
		return err
	}
	srv, err := k.serviceKeeper.Get(request, inst.ServiceHash)
	if err != nil {
		return err
	}
	return srv.RequireTaskOutputs(taskKey, outputs)
}

// Get returns the execution that matches given hash.
func (k *Keeper) Get(request cosmostypes.Request, hash hash.Hash) (*execution.Execution, error) {
	var exec *execution.Execution
	store := request.KVStore(k.storeKey)
	if !store.Has(hash) {
		return nil, fmt.Errorf("execution %q not found", hash)
	}
	return exec, codec.UnmarshalBinaryBare(store.Get(hash), &exec)
}

// List returns all executions.
func (k *Keeper) List(request cosmostypes.Request) ([]*execution.Execution, error) {
	var (
		execs []*execution.Execution
		iter  = request.KVStore(k.storeKey).Iterator(nil, nil)
	)
	for iter.Valid() {
		var exec *execution.Execution
		value := iter.Value()
		if err := codec.UnmarshalBinaryBare(value, &exec); err != nil {
			return nil, err
		}
		execs = append(execs, exec)
		iter.Next()
	}
	iter.Close()
	return execs, nil
}
