package executionsdk

import (
	"errors"
	"fmt"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/protobuf/types"
	instancesdk "github.com/mesg-foundation/engine/sdk/instance"
	runnersdk "github.com/mesg-foundation/engine/sdk/runner"
	servicesdk "github.com/mesg-foundation/engine/sdk/service"
	abci "github.com/tendermint/tendermint/abci/types"
)

const backendName = "execution"

// Backend is the execution backend.
type Backend struct {
	storeKey     *cosmostypes.KVStoreKey
	serviceBack  *servicesdk.Backend
	instanceBack *instancesdk.Backend
	runnerBack   *runnersdk.Backend
}

// NewBackend returns the backend of the execution sdk.
func NewBackend(appFactory *cosmos.AppFactory, serviceBack *servicesdk.Backend, instanceBack *instancesdk.Backend, runnerBack *runnersdk.Backend) *Backend {
	backend := &Backend{
		storeKey:     cosmostypes.NewKVStoreKey(backendName),
		serviceBack:  serviceBack,
		instanceBack: instanceBack,
		runnerBack:   runnerBack,
	}
	appBackendBasic := cosmos.NewAppModuleBasic(backendName)
	appBackend := cosmos.NewAppModule(appBackendBasic, backend.handler, backend.querier)
	appFactory.RegisterModule(appBackend)
	appFactory.RegisterStoreKey(backend.storeKey)

	return backend
}

func (s *Backend) handler(request cosmostypes.Request, msg cosmostypes.Msg) (hash.Hash, error) {
	switch msg := msg.(type) {
	case msgCreateExecution:
		exec, err := s.Create(request, msg)
		if err != nil {
			return nil, err
		}
		return exec.Hash, nil
	case msgUpdateExecution:
		exec, err := s.Update(request, msg)
		if err != nil {
			return nil, err
		}
		return exec.Hash, nil
	default:
		errmsg := fmt.Sprintf("Unrecognized execution Msg type: %v", msg.Type())
		return nil, cosmostypes.ErrUnknownRequest(errmsg)
	}
}

func (s *Backend) querier(request cosmostypes.Request, path []string, req abci.RequestQuery) (interface{}, error) {
	switch path[0] {
	case "get":
		hash, err := hash.Decode(path[1])
		if err != nil {
			return nil, err
		}
		return s.Get(request, hash)
	case "list":
		return s.List(request)
	default:
		return nil, errors.New("unknown execution query endpoint" + path[0])
	}
}

// Create creates a new execution from definition.
func (s *Backend) Create(request cosmostypes.Request, msg msgCreateExecution) (*execution.Execution, error) {
	run, err := s.runnerBack.Get(request, msg.Request.ExecutorHash)
	if err != nil {
		return nil, err
	}
	inst, err := s.instanceBack.Get(request, run.InstanceHash)
	if err != nil {
		return nil, err
	}
	srv, err := s.serviceBack.Get(request, inst.ServiceHash)
	if err != nil {
		return nil, err
	}
	// TODO: to re-implement when process is on cosmos
	// if !msg.Request.ProcessHash.IsZero() {
	// 	if _, err := s.processSDK.Get(msg.Request.ProcessHash); err != nil {
	// 		return nil, err
	// 	}
	// }
	if err := srv.RequireTaskInputs(msg.Request.TaskKey, msg.Request.Inputs); err != nil {
		return nil, err
	}
	exec := execution.New(
		msg.Request.ProcessHash,
		run.InstanceHash,
		msg.Request.ParentHash,
		msg.Request.EventHash,
		msg.Request.StepID,
		msg.Request.TaskKey,
		msg.Request.Inputs,
		msg.Request.Tags,
		msg.Request.ExecutorHash,
	)
	if err := exec.Execute(); err != nil {
		return nil, err
	}
	store := request.KVStore(s.storeKey)
	value, err := codec.MarshalBinaryBare(exec)
	if err != nil {
		return nil, err
	}
	store.Set(exec.Hash, value)
	return exec, nil
}

// Update updates a new execution from definition.
func (s *Backend) Update(request cosmostypes.Request, msg msgUpdateExecution) (*execution.Execution, error) {
	store := request.KVStore(s.storeKey)
	if !store.Has(msg.Request.Hash) {
		return nil, fmt.Errorf("execution %q doesn't exist", msg.Request.Hash)
	}
	var exec *execution.Execution
	if err := codec.UnmarshalBinaryBare(store.Get(msg.Request.Hash), &exec); err != nil {
		return nil, err
	}
	switch res := msg.Request.Result.(type) {
	case *api.UpdateExecutionRequest_Outputs:
		if err := s.validateExecutionOutput(request, exec.InstanceHash, exec.TaskKey, res.Outputs); err != nil {
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
	store.Set(exec.Hash, value)
	return exec, nil
}

func (s *Backend) validateExecutionOutput(request cosmostypes.Request, instanceHash hash.Hash, taskKey string, outputs *types.Struct) error {
	inst, err := s.instanceBack.Get(request, instanceHash)
	if err != nil {
		return err
	}
	srv, err := s.serviceBack.Get(request, inst.ServiceHash)
	if err != nil {
		return err
	}
	return srv.RequireTaskOutputs(taskKey, outputs)
}

// Get returns the execution that matches given hash.
func (s *Backend) Get(request cosmostypes.Request, hash hash.Hash) (*execution.Execution, error) {
	var exec *execution.Execution
	store := request.KVStore(s.storeKey)
	if !store.Has(hash) {
		return nil, fmt.Errorf("execution %q not found", hash)
	}
	return exec, codec.UnmarshalBinaryBare(store.Get(hash), &exec)
}

// List returns all executions.
func (s *Backend) List(request cosmostypes.Request) ([]*execution.Execution, error) {
	var (
		execs []*execution.Execution
		iter  = request.KVStore(s.storeKey).Iterator(nil, nil)
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
