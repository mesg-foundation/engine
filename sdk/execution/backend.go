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
	requestStoreKey *cosmostypes.KVStoreKey
	resultStoreKey  *cosmostypes.KVStoreKey
	serviceBack     *servicesdk.Backend
	instanceBack    *instancesdk.Backend
	runnerBack      *runnersdk.Backend
}

// NewBackend returns the backend of the execution sdk.
func NewBackend(appFactory *cosmos.AppFactory, serviceBack *servicesdk.Backend, instanceBack *instancesdk.Backend, runnerBack *runnersdk.Backend) *Backend {
	backend := &Backend{
		requestStoreKey: cosmostypes.NewKVStoreKey(backendName + ".request"),
		resultStoreKey:  cosmostypes.NewKVStoreKey(backendName + ".result"),
		serviceBack:     serviceBack,
		instanceBack:    instanceBack,
		runnerBack:      runnerBack,
	}
	appBackendBasic := cosmos.NewAppModuleBasic(backendName)
	appBackend := cosmos.NewAppModule(appBackendBasic, backend.handler, backend.querier)
	appFactory.RegisterModule(appBackend)
	appFactory.RegisterStoreKey(backend.requestStoreKey)
	appFactory.RegisterStoreKey(backend.resultStoreKey)

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
	execReq := execution.NewRequest(
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
	reqStore := request.KVStore(s.requestStoreKey)
	if reqStore.Has(execReq.Hash) {
		return nil, fmt.Errorf("execution request %q already exists", execReq.Hash)
	}
	value, err := codec.MarshalBinaryBare(execReq)
	if err != nil {
		return nil, err
	}
	reqStore.Set(execReq.Hash, value)
	return execution.ToExecution(execReq, nil), nil
}

// Update updates a new execution from definition.
func (s *Backend) Update(request cosmostypes.Request, msg msgUpdateExecution) (*execution.Execution, error) {
	reqStore := request.KVStore(s.requestStoreKey)
	if !reqStore.Has(msg.Request.Hash) {
		return nil, fmt.Errorf("execution request %q doesn't exist", msg.Request.Hash)
	}
	execRes, err := s.getExecutionResult(request.KVStore(s.resultStoreKey), msg.Request.Hash)
	if err != nil {
		return nil, err
	}
	if execRes != nil {
		return nil, fmt.Errorf("execution request %q has already a result", msg.Request.Hash)
	}

	var execReq *execution.ExecutionRequest
	if err := codec.UnmarshalBinaryBare(reqStore.Get(msg.Request.Hash), &execReq); err != nil {
		return nil, err
	}
	switch res := msg.Request.Result.(type) {
	case *api.UpdateExecutionRequest_Outputs:
		if err := s.validateExecutionOutput(request, execReq.InstanceHash, execReq.TaskKey, res.Outputs); err != nil {
			execRes = execution.NewResultWithError(execReq.Hash, err.Error())
		} else {
			execRes = execution.NewResultWithOutputs(execReq.Hash, res.Outputs)
		}
	case *api.UpdateExecutionRequest_Error:
		execRes = execution.NewResultWithError(execReq.Hash, res.Error)
	default:
		return nil, errors.New("no execution result outputs or error supplied")
	}

	value, err := codec.MarshalBinaryBare(execRes)
	if err != nil {
		return nil, err
	}
	resStore := request.KVStore(s.resultStoreKey)
	if resStore.Has(execRes.Hash) {
		return nil, fmt.Errorf("execution result %q already exists", execRes.Hash)
	}
	resStore.Set(execRes.Hash, value)
	return execution.ToExecution(execReq, execRes), nil
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
	reqStore := request.KVStore(s.requestStoreKey)
	var execReq *execution.ExecutionRequest
	if !reqStore.Has(hash) {
		return nil, fmt.Errorf("execution request %q not found", hash)
	}
	if err := codec.UnmarshalBinaryBare(reqStore.Get(hash), &execReq); err != nil {
		return nil, err
	}
	execRes, err := s.getExecutionResult(request.KVStore(s.resultStoreKey), execReq.Hash)
	if err != nil {
		return nil, err
	}
	return execution.ToExecution(execReq, execRes), nil
}

// List returns all executions.
func (s *Backend) List(request cosmostypes.Request) ([]*execution.Execution, error) {
	var (
		execs    []*execution.Execution
		iter     = request.KVStore(s.requestStoreKey).Iterator(nil, nil)
		resStore = request.KVStore(s.resultStoreKey)
	)
	for iter.Valid() {
		var execReq *execution.ExecutionRequest
		value := iter.Value()
		if err := codec.UnmarshalBinaryBare(value, &execReq); err != nil {
			return nil, err
		}
		execRes, err := s.getExecutionResult(resStore, execReq.Hash)
		if err != nil {
			return nil, err
		}
		execs = append(execs, execution.ToExecution(execReq, execRes))
		iter.Next()
	}
	iter.Close()
	return execs, nil
}

func (s *Backend) getExecutionResult(resStore cosmostypes.KVStore, requestHash hash.Hash) (*execution.ExecutionResult, error) {
	iter := resStore.Iterator(nil, nil)
	for iter.Valid() {
		var execRes *execution.ExecutionResult
		value := iter.Value()
		if err := codec.UnmarshalBinaryBare(value, &execRes); err != nil {
			return nil, err
		}
		if execRes.RequestHash.Equal(requestHash) {
			iter.Close()
			return execRes, nil
		}
		iter.Next()
	}
	iter.Close()
	return nil, nil
}
