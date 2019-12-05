package executionsdk

import (
	"errors"
	"fmt"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/codec"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/execution"
	"github.com/mesg-foundation/engine/hash"
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
	default:
		errmsg := fmt.Sprintf("Unrecognized execution msg type: %v", msg.Type())
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

// Create creates a new execution.
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
		msg.Request.ParentResultHash,
		msg.Request.EventHash,
		msg.Request.NodeKey,
		msg.Request.TaskKey,
		msg.Request.Inputs,
		msg.Request.Tags,
		msg.Request.ExecutorHash,
	)
	store := request.KVStore(s.storeKey)
	if store.Has(exec.Hash) {
		return nil, fmt.Errorf("execution %q already exists", exec.Hash)
	}
	value, err := codec.MarshalBinaryBare(exec)
	if err != nil {
		return nil, err
	}
	store.Set(exec.Hash, value)
	return exec, nil
}

// Get returns the execution that matches given hash.
func (s *Backend) Get(request cosmostypes.Request, hash hash.Hash) (*execution.Execution, error) {
	store := request.KVStore(s.storeKey)
	var exec *execution.Execution
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
